package inventory

import (
	"fmt"

	emailApp "github.com/YasserCherfaoui/darween/internal/application/email"
	companyDomain "github.com/YasserCherfaoui/darween/internal/domain/company"
	franchiseDomain "github.com/YasserCherfaoui/darween/internal/domain/franchise"
	"github.com/YasserCherfaoui/darween/internal/domain/inventory"
	productDomain "github.com/YasserCherfaoui/darween/internal/domain/product"
	userDomain "github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/pkg/errors"
)

type Service struct {
	inventoryRepo inventory.Repository
	companyRepo   companyDomain.Repository
	franchiseRepo franchiseDomain.Repository
	userRepo      userDomain.Repository
	productRepo   productDomain.Repository
	emailService  *emailApp.Service
}

func NewService(
	inventoryRepo inventory.Repository,
	companyRepo companyDomain.Repository,
	franchiseRepo franchiseDomain.Repository,
	userRepo userDomain.Repository,
	productRepo productDomain.Repository,
	emailService *emailApp.Service,
) *Service {
	return &Service{
		inventoryRepo: inventoryRepo,
		companyRepo:   companyRepo,
		franchiseRepo: franchiseRepo,
		userRepo:      userRepo,
		productRepo:   productRepo,
		emailService:  emailService,
	}
}

func (s *Service) CreateInventory(userID uint, req *CreateInventoryRequest) (*InventoryResponse, error) {
	// Validate that either company or franchise is set, but not both
	if (req.CompanyID == nil && req.FranchiseID == nil) || (req.CompanyID != nil && req.FranchiseID != nil) {
		return nil, errors.NewValidationError("must specify either company_id or franchise_id, not both")
	}

	// Check user has access to the company/franchise
	if req.CompanyID != nil {
		role, err := s.userRepo.FindUserRoleInCompany(userID, *req.CompanyID)
		if err != nil || (role.Role != userDomain.RoleOwner && role.Role != userDomain.RoleAdmin) {
			return nil, errors.NewForbiddenError("only owners and admins can create inventory")
		}
	} else if req.FranchiseID != nil {
		franchise, err := s.franchiseRepo.FindByID(*req.FranchiseID)
		if err != nil {
			return nil, errors.NewNotFoundError("franchise not found")
		}

		// Check if user is parent company admin or franchise admin
		parentRole, _ := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
		franchiseRole, err := s.userRepo.FindUserRoleInFranchise(userID, *req.FranchiseID)

		if (parentRole == nil || (parentRole.Role != userDomain.RoleOwner && parentRole.Role != userDomain.RoleAdmin)) &&
			(franchiseRole == nil || (franchiseRole.Role != userDomain.RoleOwner && franchiseRole.Role != userDomain.RoleAdmin)) {
			return nil, errors.NewForbiddenError("only owners and admins can create inventory")
		}
	}

	// Verify product variant exists
	_, err := s.productRepo.FindProductVariantByID(req.ProductVariantID)
	if err != nil {
		return nil, errors.NewNotFoundError("product variant not found")
	}

	// Create inventory
	newInventory := &inventory.Inventory{
		ProductVariantID: req.ProductVariantID,
		CompanyID:        req.CompanyID,
		FranchiseID:      req.FranchiseID,
		Stock:            req.Stock,
		ReservedStock:    0,
		IsActive:         true,
	}

	if err := s.inventoryRepo.Create(newInventory); err != nil {
		return nil, errors.NewInternalError("failed to create inventory", err)
	}

	// Log movement
	movement := &inventory.InventoryMovement{
		InventoryID:   newInventory.ID,
		MovementType:  inventory.MovementTypePurchase,
		Quantity:      req.Stock,
		PreviousStock: 0,
		NewStock:      req.Stock,
		Notes:         stringPtr("Initial stock"),
		CreatedByID:   userID,
	}
	if err := s.inventoryRepo.CreateMovement(movement); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to log inventory movement: %v\n", err)
	}

	return s.buildInventoryResponse(newInventory)
}

func (s *Service) GetInventoryByCompany(userID, companyID uint) ([]*InventoryResponse, error) {
	// Check user has access to company
	_, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return nil, errors.NewForbiddenError("you don't have access to this company")
	}

	inventories, err := s.inventoryRepo.FindByCompany(companyID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch inventory", err)
	}

	result := make([]*InventoryResponse, 0, len(inventories))
	for _, inv := range inventories {
		resp, err := s.buildInventoryResponse(inv)
		if err != nil {
			continue
		}
		result = append(result, resp)
	}

	return result, nil
}

func (s *Service) GetInventoryByFranchise(userID, franchiseID uint) ([]*InventoryResponse, error) {
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}

	// Check if user is parent company admin or franchise user
	parentRole, _ := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
	franchiseRole, _ := s.userRepo.FindUserRoleInFranchise(userID, franchiseID)

	if (parentRole == nil || (parentRole.Role != userDomain.RoleOwner && parentRole.Role != userDomain.RoleAdmin)) &&
		franchiseRole == nil {
		return nil, errors.NewForbiddenError("you don't have access to this franchise")
	}

	inventories, err := s.inventoryRepo.FindByFranchise(franchiseID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch inventory", err)
	}

	result := make([]*InventoryResponse, 0, len(inventories))
	for _, inv := range inventories {
		resp, err := s.buildInventoryResponse(inv)
		if err != nil {
			continue
		}
		result = append(result, resp)
	}

	return result, nil
}

func (s *Service) UpdateInventoryStock(userID, inventoryID uint, req *UpdateInventoryStockRequest) (*InventoryResponse, error) {
	inv, err := s.inventoryRepo.FindByID(inventoryID)
	if err != nil {
		return nil, errors.NewNotFoundError("inventory not found")
	}

	// Verify authorization
	if err := s.verifyInventoryAccess(userID, inv); err != nil {
		return nil, err
	}

	previousStock := inv.Stock
	inv.Stock = req.Stock

	if err := s.inventoryRepo.Update(inv); err != nil {
		return nil, errors.NewInternalError("failed to update inventory", err)
	}

	// Log movement
	movement := &inventory.InventoryMovement{
		InventoryID:   inventoryID,
		MovementType:  inventory.MovementTypeAdjustment,
		Quantity:      req.Stock - previousStock,
		PreviousStock: previousStock,
		NewStock:      req.Stock,
		CreatedByID:   userID,
	}
	if err := s.inventoryRepo.CreateMovement(movement); err != nil {
		fmt.Printf("Failed to log inventory movement: %v\n", err)
	}

	// Check for stock alerts
	s.checkAndSendStockAlert(inv)

	return s.buildInventoryResponse(inv)
}

func (s *Service) AdjustInventoryStock(userID, inventoryID uint, req *AdjustInventoryStockRequest) (*InventoryResponse, error) {
	inv, err := s.inventoryRepo.FindByID(inventoryID)
	if err != nil {
		return nil, errors.NewNotFoundError("inventory not found")
	}

	// Verify authorization
	if err := s.verifyInventoryAccess(userID, inv); err != nil {
		return nil, err
	}

	previousStock := inv.Stock

	// Apply adjustment
	if req.Adjustment > 0 {
		inv.AddStock(req.Adjustment)
	} else if req.Adjustment < 0 {
		amount := -req.Adjustment
		if !inv.RemoveStock(amount) {
			return nil, errors.NewValidationError("insufficient stock")
		}
	}

	if err := s.inventoryRepo.Update(inv); err != nil {
		return nil, errors.NewInternalError("failed to adjust inventory", err)
	}

	// Log movement
	movement := &inventory.InventoryMovement{
		InventoryID:   inventoryID,
		MovementType:  inventory.MovementTypeAdjustment,
		Quantity:      req.Adjustment,
		PreviousStock: previousStock,
		NewStock:      inv.Stock,
		Notes:         stringPtr(req.Notes),
		CreatedByID:   userID,
	}
	if err := s.inventoryRepo.CreateMovement(movement); err != nil {
		fmt.Printf("Failed to log inventory movement: %v\n", err)
	}

	// Check for stock alerts
	s.checkAndSendStockAlert(inv)

	return s.buildInventoryResponse(inv)
}

func (s *Service) ReserveStock(userID, inventoryID uint, req *ReserveStockRequest) (*InventoryResponse, error) {
	inv, err := s.inventoryRepo.FindByID(inventoryID)
	if err != nil {
		return nil, errors.NewNotFoundError("inventory not found")
	}

	// Verify authorization
	if err := s.verifyInventoryAccess(userID, inv); err != nil {
		return nil, err
	}

	if !inv.ReserveStock(req.Quantity) {
		return nil, errors.NewValidationError("insufficient available stock")
	}

	if err := s.inventoryRepo.Update(inv); err != nil {
		return nil, errors.NewInternalError("failed to reserve stock", err)
	}

	// Log movement
	movement := &inventory.InventoryMovement{
		InventoryID:   inventoryID,
		MovementType:  inventory.MovementTypeReserve,
		Quantity:      req.Quantity,
		PreviousStock: inv.Stock - inv.ReservedStock,
		NewStock:      inv.Stock,
		ReferenceType: stringPtr(req.ReferenceType),
		ReferenceID:   stringPtr(req.ReferenceID),
		Notes:         stringPtr(req.Notes),
		CreatedByID:   userID,
	}
	if err := s.inventoryRepo.CreateMovement(movement); err != nil {
		fmt.Printf("Failed to log inventory movement: %v\n", err)
	}

	return s.buildInventoryResponse(inv)
}

func (s *Service) ReleaseStock(userID, inventoryID uint, req *ReleaseStockRequest) (*InventoryResponse, error) {
	inv, err := s.inventoryRepo.FindByID(inventoryID)
	if err != nil {
		return nil, errors.NewNotFoundError("inventory not found")
	}

	// Verify authorization
	if err := s.verifyInventoryAccess(userID, inv); err != nil {
		return nil, err
	}

	previousStock := inv.Stock
	inv.ReleaseStock(req.Quantity)

	if err := s.inventoryRepo.Update(inv); err != nil {
		return nil, errors.NewInternalError("failed to release stock", err)
	}

	// Log movement
	movement := &inventory.InventoryMovement{
		InventoryID:   inventoryID,
		MovementType:  inventory.MovementTypeRelease,
		Quantity:      -req.Quantity,
		PreviousStock: previousStock,
		NewStock:      inv.Stock,
		Notes:         stringPtr(req.Notes),
		CreatedByID:   userID,
	}
	if err := s.inventoryRepo.CreateMovement(movement); err != nil {
		fmt.Printf("Failed to log inventory movement: %v\n", err)
	}

	return s.buildInventoryResponse(inv)
}

func (s *Service) GetInventoryMovements(inventoryID uint, limit int) ([]*InventoryMovementResponse, error) {
	movements, err := s.inventoryRepo.FindMovementsByInventory(inventoryID, limit)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch movements", err)
	}

	result := make([]*InventoryMovementResponse, 0, len(movements))
	for _, m := range movements {
		result = append(result, &InventoryMovementResponse{
			ID:            m.ID,
			InventoryID:   m.InventoryID,
			MovementType:  string(m.MovementType),
			Quantity:      m.Quantity,
			PreviousStock: m.PreviousStock,
			NewStock:      m.NewStock,
			ReferenceType: getStringValue(m.ReferenceType),
			ReferenceID:   getStringValue(m.ReferenceID),
			Notes:         getStringValue(m.Notes),
			CreatedByID:   m.CreatedByID,
			CreatedAt:     m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return result, nil
}

func (s *Service) GetInventoryMovementsWithFilters(inventoryID uint, req *MovementFilterRequest) (*InventoryMovementListResponse, error) {
	var movementType *string
	if req.MovementType != "" {
		movementType = &req.MovementType
	}

	var startDate *string
	if req.StartDate != "" {
		startDate = &req.StartDate
	}

	var endDate *string
	if req.EndDate != "" {
		endDate = &req.EndDate
	}

	movements, total, err := s.inventoryRepo.FindMovementsByInventoryWithFilters(
		inventoryID,
		movementType,
		startDate,
		endDate,
		req.Page,
		req.Limit,
	)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch movements", err)
	}

	result := make([]*InventoryMovementResponse, 0, len(movements))
	for _, m := range movements {
		result = append(result, &InventoryMovementResponse{
			ID:            m.ID,
			InventoryID:   m.InventoryID,
			MovementType:  string(m.MovementType),
			Quantity:      m.Quantity,
			PreviousStock: m.PreviousStock,
			NewStock:      m.NewStock,
			ReferenceType: getStringValue(m.ReferenceType),
			ReferenceID:   getStringValue(m.ReferenceID),
			Notes:         getStringValue(m.Notes),
			CreatedByID:   m.CreatedByID,
			CreatedAt:     m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &InventoryMovementListResponse{
		Movements:  result,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// Helper methods

func (s *Service) verifyInventoryAccess(userID uint, inv *inventory.Inventory) error {
	if inv.CompanyID != nil {
		role, err := s.userRepo.FindUserRoleInCompany(userID, *inv.CompanyID)
		if err != nil || (role.Role != userDomain.RoleOwner && role.Role != userDomain.RoleAdmin) {
			return errors.NewForbiddenError("only owners and admins can manage inventory")
		}
	} else if inv.FranchiseID != nil {
		franchise, err := s.franchiseRepo.FindByID(*inv.FranchiseID)
		if err != nil {
			return errors.NewNotFoundError("franchise not found")
		}

		parentRole, _ := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
		franchiseRole, _ := s.userRepo.FindUserRoleInFranchise(userID, *inv.FranchiseID)

		if (parentRole == nil || (parentRole.Role != userDomain.RoleOwner && parentRole.Role != userDomain.RoleAdmin)) &&
			(franchiseRole == nil || franchiseRole.Role == userDomain.RoleEmployee) {
			return errors.NewForbiddenError("you don't have access to manage this inventory")
		}
	}

	return nil
}

func (s *Service) buildInventoryResponse(inv *inventory.Inventory) (*InventoryResponse, error) {
	response := &InventoryResponse{
		ID:               inv.ID,
		ProductVariantID: inv.ProductVariantID,
		CompanyID:        inv.CompanyID,
		FranchiseID:      inv.FranchiseID,
		Stock:            inv.Stock,
		ReservedStock:    inv.ReservedStock,
		AvailableStock:   inv.GetAvailableStock(),
		IsActive:         inv.IsActive,
		CreatedAt:        inv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        inv.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Load and populate variant details from productRepo
	variant, err := s.productRepo.FindProductVariantByID(inv.ProductVariantID)
	if err == nil && variant != nil {
		response.VariantName = variant.Name
		response.VariantSKU = variant.SKU
		
		// Load product details
		product, err := s.productRepo.FindProductByID(variant.ProductID)
		if err == nil && product != nil {
			response.ProductID = product.ID
			response.ProductName = product.Name
		}
	}

	// Load franchise details if franchise inventory
	if inv.FranchiseID != nil {
		franchise, err := s.franchiseRepo.FindByID(*inv.FranchiseID)
		if err == nil && franchise != nil {
			response.FranchiseName = franchise.Name
		}
	}

	return response, nil
}

// checkAndSendStockAlert checks if stock is below threshold and sends alert email
func (s *Service) checkAndSendStockAlert(inv *inventory.Inventory) {
	if s.emailService == nil {
		return
	}

	// Default threshold: 10 units (can be made configurable per product/variant later)
	threshold := 10

	// Only send alert if stock falls below threshold
	if inv.Stock >= threshold {
		return
	}

	// Get company ID
	companyID := inv.CompanyID
	if companyID == nil && inv.FranchiseID != nil {
		// Get franchise to get company ID
		franchise, err := s.franchiseRepo.FindByID(*inv.FranchiseID)
		if err != nil {
			return
		}
		companyID = &franchise.ParentCompanyID
	}

	if companyID == nil {
		return
	}

	// Get company
	company, err := s.companyRepo.FindByID(*companyID)
	if err != nil {
		return
	}

	// Get company users (managers and admins) to send emails to
	companyUsers, err := s.userRepo.FindByCompanyID(*companyID)
	if err != nil {
		return
	}

	// Filter to managers and admins
	var recipients []string
	for _, user := range companyUsers {
		role, err := s.userRepo.FindUserRoleInCompany(user.ID, *companyID)
		if err == nil && (role.Role == userDomain.RoleOwner || role.Role == userDomain.RoleAdmin || role.Role == userDomain.RoleManager) {
			recipients = append(recipients, user.Email)
		}
	}

	if len(recipients) == 0 {
		return
	}

	// Get product and variant details
	variant, err := s.productRepo.FindProductVariantByID(inv.ProductVariantID)
	if err != nil {
		return
	}

	product, err := s.productRepo.FindProductByID(variant.ProductID)
	if err != nil {
		return
	}

	// Send stock alert email
	emailReq := &emailApp.SendStockAlertEmailRequest{
		CompanyID:     *companyID,
		To:            recipients,
		ProductName:   product.Name,
		VariantName:   variant.Name,
		CurrentStock:  inv.Stock,
		Threshold:     threshold,
		ProductDetails: map[string]interface{}{
			"product_id":   product.ID,
			"variant_id":   variant.ID,
			"product_sku":  product.SKU,
			"variant_sku":  variant.SKU,
		},
		CompanyName: company.Name,
	}

	if err := s.emailService.SendStockAlertEmail(emailReq); err != nil {
		fmt.Printf("Failed to send stock alert email: %v\n", err)
	}
}

// Helper functions

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
