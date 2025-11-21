package supplier

import (
	"fmt"
	"strings"
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/inventory"
	"github.com/YasserCherfaoui/darween/internal/domain/product"
	"github.com/YasserCherfaoui/darween/internal/domain/supplier"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"gorm.io/gorm"
)

type Service struct {
	supplierRepo  supplier.Repository
	userRepo      user.Repository
	inventoryRepo inventory.Repository
	productRepo   product.Repository
	db            *gorm.DB
}

func NewService(supplierRepo supplier.Repository, userRepo user.Repository, inventoryRepo inventory.Repository, productRepo product.Repository, db *gorm.DB) *Service {
	return &Service{
		supplierRepo:  supplierRepo,
		userRepo:      userRepo,
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
		db:            db,
	}
}

// Supplier operations
func (s *Service) CreateSupplier(userID, companyID uint, req *CreateSupplierRequest) (*SupplierResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Validate email format if provided
	if req.Email != "" && !isValidEmail(req.Email) {
		return nil, errors.NewValidationError("invalid email format")
	}

	// Create supplier
	newSupplier := req.ToSupplier(companyID)
	if !newSupplier.IsValid() {
		return nil, errors.NewValidationError("invalid supplier data")
	}

	if err := s.supplierRepo.CreateSupplier(newSupplier); err != nil {
		return nil, errors.NewInternalError("failed to create supplier", err)
	}

	return ToSupplierResponse(newSupplier, 0), nil
}

func (s *Service) GetSuppliersByCompanyID(userID, companyID uint, page, limit int) (*PaginatedResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	suppliers, total, err := s.supplierRepo.FindSuppliersByCompanyID(companyID, page, limit)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch suppliers", err)
	}

	// Convert to response DTOs
	supplierResponses := make([]*SupplierResponse, len(suppliers))
	for i, sup := range suppliers {
		// Get product count for each supplier
		products, _ := s.supplierRepo.FindProductsBySupplier(sup.ID)
		supplierResponses[i] = ToSupplierResponse(sup, len(products))
	}

	return NewPaginatedResponse(supplierResponses, total, page, limit), nil
}

func (s *Service) GetSupplierByID(userID, companyID, supplierID uint) (*SupplierResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	sup, err := s.supplierRepo.FindSupplierByIDAndCompany(supplierID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier not found")
	}

	// Get product count
	products, _ := s.supplierRepo.FindProductsBySupplier(sup.ID)
	return ToSupplierResponse(sup, len(products)), nil
}

func (s *Service) UpdateSupplier(userID, companyID, supplierID uint, req *UpdateSupplierRequest) (*SupplierResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Get existing supplier
	existingSupplier, err := s.supplierRepo.FindSupplierByIDAndCompany(supplierID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier not found")
	}

	// Update fields
	if req.Name != "" {
		existingSupplier.Name = req.Name
	}
	if req.ContactPerson != "" {
		existingSupplier.ContactPerson = req.ContactPerson
	}
	if req.Email != "" {
		if !isValidEmail(req.Email) {
			return nil, errors.NewValidationError("invalid email format")
		}
		existingSupplier.Email = req.Email
	}
	if req.Phone != "" {
		existingSupplier.Phone = req.Phone
	}
	if req.Address != "" {
		existingSupplier.Address = req.Address
	}
	if req.IsActive != nil {
		existingSupplier.IsActive = *req.IsActive
	}

	if err := s.supplierRepo.UpdateSupplier(existingSupplier); err != nil {
		return nil, errors.NewInternalError("failed to update supplier", err)
	}

	// Get product count
	products, _ := s.supplierRepo.FindProductsBySupplier(existingSupplier.ID)
	return ToSupplierResponse(existingSupplier, len(products)), nil
}

func (s *Service) DeleteSupplier(userID, companyID, supplierID uint) error {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return err
	}

	// Check if supplier exists
	_, err := s.supplierRepo.FindSupplierByIDAndCompany(supplierID, companyID)
	if err != nil {
		return errors.NewNotFoundError("supplier not found")
	}

	if err := s.supplierRepo.SoftDeleteSupplier(supplierID); err != nil {
		return errors.NewInternalError("failed to delete supplier", err)
	}

	return nil
}

func (s *Service) GetSupplierProducts(userID, companyID, supplierID uint) (*SupplierWithProductsResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// Get supplier
	sup, err := s.supplierRepo.FindSupplierByIDAndCompany(supplierID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier not found")
	}

	// Get products for this supplier
	products, err := s.supplierRepo.FindProductsBySupplier(supplierID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch supplier products", err)
	}

	// Convert to response
	productInfos := make([]SupplierProductInfo, len(products))
	for i, p := range products {
		productInfos[i] = SupplierProductInfo{
			ID:           p.ID,
			Name:         p.Name,
			SKU:          p.SKU,
			SupplierCost: p.SupplierCost,
		}
	}

	return &SupplierWithProductsResponse{
		ID:            sup.ID,
		CompanyID:     sup.CompanyID,
		Name:          sup.Name,
		ContactPerson: sup.ContactPerson,
		Email:         sup.Email,
		Phone:         sup.Phone,
		Address:       sup.Address,
		IsActive:      sup.IsActive,
		Products:      productInfos,
	}, nil
}

// Helper function to check user access to company
func (s *Service) checkUserCompanyAccess(userID, companyID uint, minRole user.Role) error {
	userRole, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return errors.NewForbiddenError("you don't have access to this company")
	}

	// Check if user has sufficient role
	if !s.hasSufficientRole(userRole.Role, minRole) {
		return errors.NewForbiddenError("insufficient permissions for this operation")
	}

	return nil
}

// Helper function to check role hierarchy
func (s *Service) hasSufficientRole(userRole, requiredRole user.Role) bool {
	roleHierarchy := map[user.Role]int{
		user.RoleEmployee: 1,
		user.RoleManager:  2,
		user.RoleAdmin:    3,
		user.RoleOwner:    4,
	}

	return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	// Basic email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// Helper function to create inventory movement
func (s *Service) createInventoryMovement(tx *gorm.DB, inventoryID uint, movementType inventory.MovementType, quantity int, previousStock, newStock int, referenceType, referenceID string, createdByID uint) error {
	movement := &inventory.InventoryMovement{
		InventoryID:   inventoryID,
		MovementType:  movementType,
		Quantity:      quantity,
		PreviousStock: previousStock,
		NewStock:      newStock,
		ReferenceType: &referenceType,
		ReferenceID:   &referenceID,
		CreatedByID:   createdByID,
	}
	return tx.Create(movement).Error
}

// Helper function to generate bill number
func generateBillNumber(companyID, billID uint) string {
	timestamp := time.Now().Format("20060102")
	return fmt.Sprintf("BILL-%d-%s-%d", companyID, timestamp, billID)
}

// Helper function for string pointer
func stringPtr(s string) *string {
	return &s
}

// SupplierBill operations

func (s *Service) CreateSupplierBill(userID, companyID uint, req *CreateSupplierBillRequest) (*SupplierBillResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Validate supplier
	_, err := s.supplierRepo.FindSupplierByIDAndCompany(req.SupplierID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier not found")
	}

	// Validate items and product variants
	for _, itemReq := range req.Items {
		variant, err := s.productRepo.FindProductVariantByID(itemReq.ProductVariantID)
		if err != nil {
			return nil, errors.NewNotFoundError(fmt.Sprintf("product variant %d not found", itemReq.ProductVariantID))
		}
		// Verify variant belongs to company (optional check)
		_ = variant // Use variant if needed
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create bill items
	billItems := make([]*supplier.SupplierBillItem, len(req.Items))
	totalAmount := 0.0
	for i, itemReq := range req.Items {
		item := &supplier.SupplierBillItem{
			ProductVariantID: itemReq.ProductVariantID,
			Quantity:         itemReq.Quantity,
			UnitCost:         itemReq.UnitCost,
		}
		item.CalculateTotal()
		billItems[i] = item
		totalAmount += item.TotalCost
	}

	// Create bill with temporary bill number (will be updated after getting ID)
	// Use timestamp-based temp value to ensure uniqueness
	tempBillNumber := fmt.Sprintf("BILL-TEMP-%d-%d", companyID, time.Now().UnixNano())
	bill := &supplier.SupplierBill{
		CompanyID:     companyID,
		SupplierID:    req.SupplierID,
		BillNumber:    tempBillNumber, // Temporary value, will be updated
		TotalAmount:   totalAmount,
		PaidAmount:    req.PaidAmount,
		PaymentStatus: supplier.PaymentStatusUnpaid,
		BillStatus:    supplier.BillStatusDraft,
		Notes:         req.Notes,
		CreatedByID:   userID,
	}
	bill.CalculatePendingAmount()
	bill.UpdatePaymentStatus()

	if !bill.IsValid() {
		tx.Rollback()
		return nil, errors.NewValidationError("invalid bill data")
	}

	// Create bill first to get ID
	if err := tx.Create(bill).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to create bill", err)
	}

	// Update bill number with actual ID
	bill.BillNumber = generateBillNumber(companyID, bill.ID)
	if err := tx.Save(bill).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update bill number", err)
	}

	// Create bill items and update inventory
	for _, item := range billItems {
		item.SupplierBillID = bill.ID
		if err := tx.Create(item).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to create bill item", err)
		}

		// Update inventory (purchase adds stock)
		inv, err := s.inventoryRepo.FindByVariantAndCompany(item.ProductVariantID, companyID)
		if err != nil {
			// If inventory doesn't exist, create it
			inv = &inventory.Inventory{
				ProductVariantID: item.ProductVariantID,
				CompanyID:        &companyID,
				Stock:            0,
				ReservedStock:    0,
				IsActive:         true,
			}
			if err := tx.Create(inv).Error; err != nil {
				tx.Rollback()
				return nil, errors.NewInternalError("failed to create inventory", err)
			}
		}

		previousStock := inv.Stock
		inv.AddStock(item.Quantity)
		if err := tx.Save(inv).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to update inventory", err)
		}

		// Create inventory movement
		billIDStr := fmt.Sprintf("%d", bill.ID)
		if err := s.createInventoryMovement(tx, inv.ID, inventory.MovementTypePurchase, item.Quantity, previousStock, inv.Stock, "supplier_bill", billIDStr, userID); err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to create inventory movement", err)
		}
	}

	// Mark bill as completed if initial payment is provided
	if req.PaidAmount > 0 {
		bill.AddPayment(req.PaidAmount)
		bill.Complete()
	} else {
		bill.Complete()
	}

	if err := tx.Save(bill).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update bill status", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.NewInternalError("failed to commit transaction", err)
	}

	// Reload bill with items and supplier
	createdBill, err := s.supplierRepo.FindSupplierBillByIDAndCompany(bill.ID, companyID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch created bill", err)
	}

	return ToSupplierBillResponse(createdBill), nil
}

func (s *Service) GetSupplierBillByID(userID, companyID, billID uint) (*SupplierBillResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	bill, err := s.supplierRepo.FindSupplierBillByIDAndCompany(billID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier bill not found")
	}

	return ToSupplierBillResponse(bill), nil
}

func (s *Service) ListSupplierBills(userID, companyID, supplierID uint, page, limit int) (*PaginatedResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	var bills []*supplier.SupplierBill
	var total int64
	var err error

	if supplierID > 0 {
		bills, total, err = s.supplierRepo.FindSupplierBillsBySupplier(supplierID, companyID, page, limit)
	} else {
		bills, total, err = s.supplierRepo.FindSupplierBillsByCompany(companyID, page, limit)
	}

	if err != nil {
		return nil, errors.NewInternalError("failed to fetch supplier bills", err)
	}

	billResponses := make([]*SupplierBillResponse, len(bills))
	for i, bill := range bills {
		billResponses[i] = ToSupplierBillResponse(bill)
	}

	return NewPaginatedResponse(billResponses, total, page, limit), nil
}

func (s *Service) UpdateSupplierBill(userID, companyID, billID uint, req *UpdateSupplierBillRequest) (*SupplierBillResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Get existing bill
	existingBill, err := s.supplierRepo.FindSupplierBillByIDAndCompany(billID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier bill not found")
	}

	// Only allow editing draft bills
	if existingBill.BillStatus != supplier.BillStatusDraft {
		return nil, errors.NewValidationError("can only edit draft bills")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// If items are provided, update them
	if req.Items != nil && len(req.Items) > 0 {
		// Get existing items
		existingItems, err := s.supplierRepo.FindSupplierBillItemsByBill(billID)
		if err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to fetch existing items", err)
		}

		// Reverse inventory for existing items
		for _, existingItem := range existingItems {
			inv, err := s.inventoryRepo.FindByVariantAndCompany(existingItem.ProductVariantID, companyID)
			if err == nil {
				previousStock := inv.Stock
				if !inv.RemoveStock(existingItem.Quantity) {
					tx.Rollback()
					return nil, errors.NewValidationError(fmt.Sprintf("insufficient stock to reverse for variant %d", existingItem.ProductVariantID))
				}
				if err := tx.Save(inv).Error; err != nil {
					tx.Rollback()
					return nil, errors.NewInternalError("failed to update inventory", err)
				}

				// Create reversal movement
				billIDStr := fmt.Sprintf("%d", billID)
				if err := s.createInventoryMovement(tx, inv.ID, inventory.MovementTypeAdjustment, -existingItem.Quantity, previousStock, inv.Stock, "supplier_bill_update", billIDStr, userID); err != nil {
					tx.Rollback()
					return nil, errors.NewInternalError("failed to create inventory movement", err)
				}
			}

			// Delete existing item
			if err := tx.Delete(existingItem).Error; err != nil {
				tx.Rollback()
				return nil, errors.NewInternalError("failed to delete bill item", err)
			}
		}

		// Create new items
		newTotalAmount := 0.0
		for _, itemReq := range req.Items {
			// Validate variant
			_, err := s.productRepo.FindProductVariantByID(itemReq.ProductVariantID)
			if err != nil {
				tx.Rollback()
				return nil, errors.NewNotFoundError(fmt.Sprintf("product variant %d not found", itemReq.ProductVariantID))
			}

			item := &supplier.SupplierBillItem{
				SupplierBillID:   billID,
				ProductVariantID: itemReq.ProductVariantID,
				Quantity:         itemReq.Quantity,
				UnitCost:         itemReq.UnitCost,
			}
			item.CalculateTotal()
			newTotalAmount += item.TotalCost

			if err := tx.Create(item).Error; err != nil {
				tx.Rollback()
				return nil, errors.NewInternalError("failed to create bill item", err)
			}

			// Update inventory
			inv, err := s.inventoryRepo.FindByVariantAndCompany(item.ProductVariantID, companyID)
			if err != nil {
				// Create inventory if doesn't exist
				inv = &inventory.Inventory{
					ProductVariantID: item.ProductVariantID,
					CompanyID:        &companyID,
					Stock:            0,
					ReservedStock:    0,
					IsActive:         true,
				}
				if err := tx.Create(inv).Error; err != nil {
					tx.Rollback()
					return nil, errors.NewInternalError("failed to create inventory", err)
				}
			}

			previousStock := inv.Stock
			inv.AddStock(item.Quantity)
			if err := tx.Save(inv).Error; err != nil {
				tx.Rollback()
				return nil, errors.NewInternalError("failed to update inventory", err)
			}

			// Create inventory movement
			billIDStr := fmt.Sprintf("%d", billID)
			if err := s.createInventoryMovement(tx, inv.ID, inventory.MovementTypePurchase, item.Quantity, previousStock, inv.Stock, "supplier_bill_update", billIDStr, userID); err != nil {
				tx.Rollback()
				return nil, errors.NewInternalError("failed to create inventory movement", err)
			}
		}

		// Update bill total
		existingBill.TotalAmount = newTotalAmount
		existingBill.CalculatePendingAmount()
		existingBill.UpdatePaymentStatus()
	}

	// Update notes if provided
	if req.Notes != "" {
		existingBill.Notes = req.Notes
	}

	// Update bill status if provided
	if req.BillStatus != nil {
		status := supplier.BillStatus(*req.BillStatus)
		if !status.IsValid() {
			tx.Rollback()
			return nil, errors.NewValidationError("invalid bill status")
		}
		existingBill.BillStatus = status
	}

	if err := tx.Save(existingBill).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update bill", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.NewInternalError("failed to commit transaction", err)
	}

	// Reload bill
	updatedBill, err := s.supplierRepo.FindSupplierBillByIDAndCompany(billID, companyID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch updated bill", err)
	}

	return ToSupplierBillResponse(updatedBill), nil
}

func (s *Service) DeleteSupplierBill(userID, companyID, billID uint) error {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return err
	}

	// Get existing bill
	existingBill, err := s.supplierRepo.FindSupplierBillByIDAndCompany(billID, companyID)
	if err != nil {
		return errors.NewNotFoundError("supplier bill not found")
	}

	// Only allow deleting draft bills or if no payment has been made
	if existingBill.BillStatus != supplier.BillStatusDraft && existingBill.PaidAmount > 0 {
		return errors.NewValidationError("cannot delete bill that has payments. Cancel it instead.")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get bill items
	items, err := s.supplierRepo.FindSupplierBillItemsByBill(billID)
	if err != nil {
		tx.Rollback()
		return errors.NewInternalError("failed to fetch bill items", err)
	}

	// Restore inventory for all items
	for _, item := range items {
		inv, err := s.inventoryRepo.FindByVariantAndCompany(item.ProductVariantID, companyID)
		if err == nil {
			previousStock := inv.Stock
			if !inv.RemoveStock(item.Quantity) {
				// If we can't remove stock (e.g., already sold), still record the movement
				// but don't fail the deletion
			} else {
				if err := tx.Save(inv).Error; err != nil {
					tx.Rollback()
					return errors.NewInternalError("failed to restore inventory", err)
				}
			}

			// Create deletion movement
			billIDStr := fmt.Sprintf("%d", billID)
			if err := s.createInventoryMovement(tx, inv.ID, inventory.MovementTypeAdjustment, -item.Quantity, previousStock, inv.Stock, "supplier_bill_delete", billIDStr, userID); err != nil {
				tx.Rollback()
				return errors.NewInternalError("failed to create inventory movement", err)
			}
		}
	}

	// Delete bill (items will be deleted by CASCADE)
	if err := tx.Delete(existingBill).Error; err != nil {
		tx.Rollback()
		return errors.NewInternalError("failed to delete bill", err)
	}

	if err := tx.Commit().Error; err != nil {
		return errors.NewInternalError("failed to commit transaction", err)
	}

	return nil
}

func (s *Service) AddBillItem(userID, companyID, billID uint, req *SupplierBillItemRequest) (*SupplierBillItemResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Get existing bill
	bill, err := s.supplierRepo.FindSupplierBillByIDAndCompany(billID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier bill not found")
	}

	if bill.BillStatus != supplier.BillStatusDraft {
		return nil, errors.NewValidationError("can only add items to draft bills")
	}

	// Validate variant
	_, err = s.productRepo.FindProductVariantByID(req.ProductVariantID)
	if err != nil {
		return nil, errors.NewNotFoundError("product variant not found")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create bill item
	item := &supplier.SupplierBillItem{
		SupplierBillID:   billID,
		ProductVariantID: req.ProductVariantID,
		Quantity:         req.Quantity,
		UnitCost:         req.UnitCost,
	}
	item.CalculateTotal()

	if err := tx.Create(item).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to create bill item", err)
	}

	// Update inventory
	inv, err := s.inventoryRepo.FindByVariantAndCompany(req.ProductVariantID, companyID)
	if err != nil {
		// Create inventory if doesn't exist
		inv = &inventory.Inventory{
			ProductVariantID: req.ProductVariantID,
			CompanyID:        &companyID,
			Stock:            0,
			ReservedStock:    0,
			IsActive:         true,
		}
		if err := tx.Create(inv).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to create inventory", err)
		}
	}

	previousStock := inv.Stock
	inv.AddStock(req.Quantity)
	if err := tx.Save(inv).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update inventory", err)
	}

	// Create inventory movement
	billIDStr := fmt.Sprintf("%d", billID)
	if err := s.createInventoryMovement(tx, inv.ID, inventory.MovementTypePurchase, req.Quantity, previousStock, inv.Stock, "supplier_bill", billIDStr, userID); err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to create inventory movement", err)
	}

	// Update bill total
	bill.TotalAmount += item.TotalCost
	bill.CalculatePendingAmount()
	bill.UpdatePaymentStatus()

	if err := tx.Save(bill).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update bill", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.NewInternalError("failed to commit transaction", err)
	}

	return ToSupplierBillItemResponse(item), nil
}

func (s *Service) UpdateBillItem(userID, companyID, billID, itemID uint, req *SupplierBillItemRequest) (*SupplierBillItemResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Get existing bill
	bill, err := s.supplierRepo.FindSupplierBillByIDAndCompany(billID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier bill not found")
	}

	if bill.BillStatus != supplier.BillStatusDraft {
		return nil, errors.NewValidationError("can only update items in draft bills")
	}

	// Get existing item
	existingItem, err := s.supplierRepo.FindSupplierBillItemByID(itemID)
	if err != nil {
		return nil, errors.NewNotFoundError("bill item not found")
	}

	if existingItem.SupplierBillID != billID {
		return nil, errors.NewForbiddenError("item does not belong to this bill")
	}

	// Validate variant
	_, err = s.productRepo.FindProductVariantByID(req.ProductVariantID)
	if err != nil {
		return nil, errors.NewNotFoundError("product variant not found")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Reverse old inventory
	oldInv, err := s.inventoryRepo.FindByVariantAndCompany(existingItem.ProductVariantID, companyID)
	if err == nil {
		previousStock := oldInv.Stock
		if !oldInv.RemoveStock(existingItem.Quantity) {
			tx.Rollback()
			return nil, errors.NewValidationError("insufficient stock to reverse")
		}
		if err := tx.Save(oldInv).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to update inventory", err)
		}

		// Create reversal movement
		billIDStr := fmt.Sprintf("%d", billID)
		if err := s.createInventoryMovement(tx, oldInv.ID, inventory.MovementTypeAdjustment, -existingItem.Quantity, previousStock, oldInv.Stock, "supplier_bill_update", billIDStr, userID); err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to create inventory movement", err)
		}
	}

	// Update item
	oldTotal := existingItem.TotalCost
	existingItem.ProductVariantID = req.ProductVariantID
	existingItem.Quantity = req.Quantity
	existingItem.UnitCost = req.UnitCost
	existingItem.CalculateTotal()

	if err := tx.Save(existingItem).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update bill item", err)
	}

	// Update new inventory
	newInv, err := s.inventoryRepo.FindByVariantAndCompany(req.ProductVariantID, companyID)
	if err != nil {
		// Create inventory if doesn't exist
		newInv = &inventory.Inventory{
			ProductVariantID: req.ProductVariantID,
			CompanyID:        &companyID,
			Stock:            0,
			ReservedStock:    0,
			IsActive:         true,
		}
		if err := tx.Create(newInv).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to create inventory", err)
		}
	}

	previousStock := newInv.Stock
	newInv.AddStock(req.Quantity)
	if err := tx.Save(newInv).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update inventory", err)
	}

	// Create inventory movement
	billIDStr := fmt.Sprintf("%d", billID)
	if err := s.createInventoryMovement(tx, newInv.ID, inventory.MovementTypePurchase, req.Quantity, previousStock, newInv.Stock, "supplier_bill_update", billIDStr, userID); err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to create inventory movement", err)
	}

	// Update bill total
	bill.TotalAmount = bill.TotalAmount - oldTotal + existingItem.TotalCost
	bill.CalculatePendingAmount()
	bill.UpdatePaymentStatus()

	if err := tx.Save(bill).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update bill", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.NewInternalError("failed to commit transaction", err)
	}

	return ToSupplierBillItemResponse(existingItem), nil
}

func (s *Service) RemoveBillItem(userID, companyID, billID, itemID uint) error {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return err
	}

	// Get existing bill
	bill, err := s.supplierRepo.FindSupplierBillByIDAndCompany(billID, companyID)
	if err != nil {
		return errors.NewNotFoundError("supplier bill not found")
	}

	if bill.BillStatus != supplier.BillStatusDraft {
		return errors.NewValidationError("can only remove items from draft bills")
	}

	// Get existing item
	existingItem, err := s.supplierRepo.FindSupplierBillItemByID(itemID)
	if err != nil {
		return errors.NewNotFoundError("bill item not found")
	}

	if existingItem.SupplierBillID != billID {
		return errors.NewForbiddenError("item does not belong to this bill")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Reverse inventory
	inv, err := s.inventoryRepo.FindByVariantAndCompany(existingItem.ProductVariantID, companyID)
	if err == nil {
		previousStock := inv.Stock
		if !inv.RemoveStock(existingItem.Quantity) {
			// If we can't remove stock, still record the movement but don't fail
		} else {
			if err := tx.Save(inv).Error; err != nil {
				tx.Rollback()
				return errors.NewInternalError("failed to restore inventory", err)
			}
		}

		// Create reversal movement
		billIDStr := fmt.Sprintf("%d", billID)
		if err := s.createInventoryMovement(tx, inv.ID, inventory.MovementTypeAdjustment, -existingItem.Quantity, previousStock, inv.Stock, "supplier_bill_update", billIDStr, userID); err != nil {
			tx.Rollback()
			return errors.NewInternalError("failed to create inventory movement", err)
		}
	}

	// Delete item
	if err := tx.Delete(existingItem).Error; err != nil {
		tx.Rollback()
		return errors.NewInternalError("failed to delete bill item", err)
	}

	// Update bill total
	bill.TotalAmount -= existingItem.TotalCost
	bill.CalculatePendingAmount()
	bill.UpdatePaymentStatus()

	if err := tx.Save(bill).Error; err != nil {
		tx.Rollback()
		return errors.NewInternalError("failed to update bill", err)
	}

	if err := tx.Commit().Error; err != nil {
		return errors.NewInternalError("failed to commit transaction", err)
	}

	return nil
}

func (s *Service) RecordSupplierPayment(userID, companyID, supplierID uint, req *RecordSupplierPaymentRequest) (*SupplierPaymentResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Validate supplier
	_, err := s.supplierRepo.FindSupplierByIDAndCompany(supplierID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier not found")
	}

	if supplierID != req.SupplierID {
		return nil, errors.NewValidationError("supplier ID mismatch")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create payment
	payment := &supplier.SupplierPayment{
		SupplierID:    req.SupplierID,
		CompanyID:     companyID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: supplier.PaymentTransactionStatusCompleted,
		Reference:     req.Reference,
		Notes:         req.Notes,
		CreatedByID:   userID,
	}

	if !payment.IsValid() {
		tx.Rollback()
		return nil, errors.NewValidationError("invalid payment data")
	}

	if err := tx.Create(payment).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to create payment", err)
	}

	// Get unpaid bills for this supplier (oldest first for FIFO distribution)
	unpaidBills, err := s.supplierRepo.FindUnpaidBillsBySupplier(supplierID, companyID)
	if err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to fetch unpaid bills", err)
	}

	// Distribute payment to bills (FIFO: oldest bills first)
	remainingAmount := req.Amount
	distributions := []*supplier.SupplierPaymentDistribution{}

	for _, bill := range unpaidBills {
		if remainingAmount <= 0 {
			break
		}

		pendingAmount := bill.TotalAmount - bill.PaidAmount
		if pendingAmount <= 0 {
			continue
		}

		distributionAmount := remainingAmount
		if distributionAmount > pendingAmount {
			distributionAmount = pendingAmount
		}

		// Create distribution
		distribution := &supplier.SupplierPaymentDistribution{
			SupplierPaymentID: payment.ID,
			SupplierBillID:    bill.ID,
			Amount:            distributionAmount,
		}

		if err := tx.Create(distribution).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to create payment distribution", err)
		}

		distributions = append(distributions, distribution)

		// Update bill paid amount
		bill.AddPayment(distributionAmount)
		if err := tx.Save(bill).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to update bill", err)
		}

		remainingAmount -= distributionAmount
	}

	if remainingAmount > 0 {
		// Payment exceeds all pending bills - this is acceptable, but log it
		// Could store the excess amount or return a warning
	}

	payment.MarkCompleted()
	if err := tx.Save(payment).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update payment", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.NewInternalError("failed to commit transaction", err)
	}

	// Reload payment with distributions
	createdPayment, err := s.supplierRepo.FindSupplierPaymentByID(payment.ID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch created payment", err)
	}

	return ToSupplierPaymentResponse(createdPayment), nil
}

func (s *Service) GetSupplierOutstandingBalance(userID, companyID, supplierID uint) (*SupplierOutstandingBalanceResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// Validate supplier
	_, err := s.supplierRepo.FindSupplierByIDAndCompany(supplierID, companyID)
	if err != nil {
		return nil, errors.NewNotFoundError("supplier not found")
	}

	// Calculate outstanding balance
	outstandingAmount, err := s.supplierRepo.CalculateSupplierOutstandingBalance(supplierID, companyID)
	if err != nil {
		return nil, errors.NewInternalError("failed to calculate outstanding balance", err)
	}

	return &SupplierOutstandingBalanceResponse{
		SupplierID:        supplierID,
		CompanyID:         companyID,
		OutstandingAmount: outstandingAmount,
	}, nil
}

