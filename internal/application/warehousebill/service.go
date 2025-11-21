package warehousebill

import (
	"encoding/json"
	"fmt"
	"time"

	emailApp "github.com/YasserCherfaoui/darween/internal/application/email"
	companyDomain "github.com/YasserCherfaoui/darween/internal/domain/company"
	franchiseDomain "github.com/YasserCherfaoui/darween/internal/domain/franchise"
	"github.com/YasserCherfaoui/darween/internal/domain/inventory"
	productDomain "github.com/YasserCherfaoui/darween/internal/domain/product"
	userDomain "github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/internal/domain/warehousebill"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"gorm.io/gorm"
)

type Service struct {
	warehouseBillRepo warehousebill.Repository
	inventoryRepo     inventory.Repository
	companyRepo       companyDomain.Repository
	franchiseRepo     franchiseDomain.Repository
	userRepo          userDomain.Repository
	productRepo       productDomain.Repository
	emailService      *emailApp.Service
	db                *gorm.DB
}

func NewService(
	warehouseBillRepo warehousebill.Repository,
	inventoryRepo inventory.Repository,
	companyRepo companyDomain.Repository,
	franchiseRepo franchiseDomain.Repository,
	userRepo userDomain.Repository,
	productRepo productDomain.Repository,
	emailService *emailApp.Service,
	db *gorm.DB,
) *Service {
	return &Service{
		warehouseBillRepo: warehouseBillRepo,
		inventoryRepo:     inventoryRepo,
		companyRepo:       companyRepo,
		franchiseRepo:     franchiseRepo,
		userRepo:          userRepo,
		productRepo:       productRepo,
		emailService:      emailService,
		db:                db,
	}
}

// CreateExitBill creates an exit bill from warehouse to franchise
func (s *Service) CreateExitBill(userID, companyID uint, req *CreateExitBillRequest) (*WarehouseBillResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, userDomain.RoleManager); err != nil {
		return nil, err
	}

	// Validate franchise belongs to company
	franchise, err := s.franchiseRepo.FindByID(req.FranchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}
	if franchise.ParentCompanyID != companyID {
		return nil, errors.NewForbiddenError("franchise does not belong to this company")
	}

	// Validate items and product variants, check inventory availability
	// Collect all validation issues instead of returning on first error
	var validationIssues []ValidationIssue
	
	for i, itemReq := range req.Items {
		variant, err := s.productRepo.FindProductVariantByID(itemReq.ProductVariantID)
		if err != nil {
			validationIssues = append(validationIssues, ValidationIssue{
				ItemIndex: i,
				VariantID: itemReq.ProductVariantID,
				Message:   fmt.Sprintf("Product variant %d not found", itemReq.ProductVariantID),
			})
			continue
		}

		// Get product information for better error messages
		product, err := s.productRepo.FindProductByID(variant.ProductID)
		if err != nil {
			validationIssues = append(validationIssues, ValidationIssue{
				ItemIndex:  i,
				VariantID:  variant.ID,
				VariantSKU: variant.SKU,
				Message:    fmt.Sprintf("Product not found for variant SKU: %s", variant.SKU),
			})
			continue
		}

		// Check company inventory availability
		companyInv, err := s.inventoryRepo.FindByVariantAndCompany(itemReq.ProductVariantID, companyID)
		if err != nil {
			// Inventory not found means available quantity is 0
			validationIssues = append(validationIssues, ValidationIssue{
				ItemIndex:   i,
				VariantID:   variant.ID,
				VariantSKU:  variant.SKU,
				ProductSKU:  product.SKU,
				ProductName: product.Name,
				AvailableQty: 0,
				RequiredQty: itemReq.Quantity,
				Message:     fmt.Sprintf("Insufficient inventory for product '%s' (SKU: %s), variant SKU: %s. Available quantity: 0, Required: %d", product.Name, product.SKU, variant.SKU, itemReq.Quantity),
			})
			continue
		}
		
		if companyInv.GetAvailableStock() < itemReq.Quantity {
			validationIssues = append(validationIssues, ValidationIssue{
				ItemIndex:   i,
				VariantID:   variant.ID,
				VariantSKU:  variant.SKU,
				ProductSKU:  product.SKU,
				ProductName: product.Name,
				AvailableQty: companyInv.GetAvailableStock(),
				RequiredQty: itemReq.Quantity,
				Message:     fmt.Sprintf("Insufficient stock for product '%s' (SKU: %s), variant SKU: %s. Available quantity: %d, Required: %d", product.Name, product.SKU, variant.SKU, companyInv.GetAvailableStock(), itemReq.Quantity),
			})
		}
	}

	// If there are validation issues, return them all
	if len(validationIssues) > 0 {
		// Serialize issues to JSON for error message
		issuesJSON, err := json.Marshal(ValidationErrorsResponse{Issues: validationIssues})
		if err != nil {
			return nil, errors.NewInternalError("failed to serialize validation errors", err)
		}
		return nil, errors.NewValidationErrorsError(string(issuesJSON))
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create bill items
	billItems := make([]warehousebill.WarehouseBillItem, len(req.Items))
	totalAmount := 0.0
	for i, itemReq := range req.Items {
		item := warehousebill.WarehouseBillItem{
			ProductVariantID: itemReq.ProductVariantID,
			Quantity:         itemReq.Quantity,
			UnitPrice:        itemReq.UnitPrice,
			DiscrepancyType:  warehousebill.DiscrepancyTypeNone,
		}
		item.CalculateTotal()
		billItems[i] = item
		totalAmount += item.TotalAmount
	}

	// Create bill
	bill := &warehousebill.WarehouseBill{
		CompanyID:          companyID,
		FranchiseID:        req.FranchiseID,
		BillType:           warehousebill.BillTypeExit,
		Status:             warehousebill.BillStatusDraft,
		VerificationStatus: warehousebill.VerificationStatusPending,
		TotalAmount:        totalAmount,
		Notes:              req.Notes,
		CreatedByID:        userID,
		Items:              billItems,
	}

	if !bill.IsValid() {
		tx.Rollback()
		return nil, errors.NewValidationError("invalid bill data")
	}

	if err := s.warehouseBillRepo.Create(bill); err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to create exit bill", err)
	}

	// Reserve stock and create inventory movements for draft bill
	for _, item := range bill.Items {
		// Get inventory within transaction (we already validated it exists)
		var companyInv inventory.Inventory
		if err := tx.Where("product_variant_id = ? AND company_id = ?", item.ProductVariantID, companyID).First(&companyInv).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError(fmt.Sprintf("failed to fetch inventory for variant %d", item.ProductVariantID), err)
		}

		// Reserve stock for draft bill
		previousReservedStock := companyInv.ReservedStock
		if !companyInv.ReserveStock(item.Quantity) {
			tx.Rollback()
			return nil, errors.NewInternalError(fmt.Sprintf("failed to reserve stock for variant %d", item.ProductVariantID), err)
		}

		// Update inventory in transaction
		if err := tx.Save(&companyInv).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to update inventory", err)
		}

		// Create inventory movement to track reservation
		refType := "warehouse_bill"
		refID := fmt.Sprintf("%d", bill.ID)
		notes := fmt.Sprintf("Stock reserved for draft exit bill (Bill ID: %d)", bill.ID)
		movement := &inventory.InventoryMovement{
			InventoryID:   companyInv.ID,
			MovementType:  inventory.MovementTypeReserve,
			Quantity:      item.Quantity,
			PreviousStock: previousReservedStock,
			NewStock:      companyInv.ReservedStock,
			ReferenceType: &refType,
			ReferenceID:   &refID,
			Notes:         &notes,
			CreatedByID:   userID,
		}
		if err := tx.Create(movement).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to create inventory movement", err)
		}
	}

	tx.Commit()
	
	// Enrich response with product/variant details
	response := ToWarehouseBillResponse(bill)
	if err := s.enrichWarehouseBillResponse(response); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to enrich bill response: %v\n", err)
	}
	
	return response, nil
}

// CreateEntryBill creates an entry bill linked to an exit bill
func (s *Service) CreateEntryBill(userID, franchiseID uint, req *CreateEntryBillRequest) (*WarehouseBillResponse, error) {
	// Check user authorization (franchise access)
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}

	franchiseRole, err := s.userRepo.FindUserRoleInFranchise(userID, franchiseID)
	if err != nil || franchiseRole == nil {
		// Check parent company access
		parentRole, err := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
		if err != nil || parentRole == nil {
			return nil, errors.NewForbiddenError("you don't have access to this franchise")
		}
	}

	// Find exit bill
	exitBill, err := s.warehouseBillRepo.FindByID(req.ExitBillID)
	if err != nil {
		return nil, errors.NewNotFoundError("exit bill not found")
	}

	// Validate exit bill
	if !exitBill.IsExitBill() {
		return nil, errors.NewValidationError("specified bill is not an exit bill")
	}
	if exitBill.FranchiseID != franchiseID {
		return nil, errors.NewForbiddenError("exit bill does not belong to this franchise")
	}
	if exitBill.Status != warehousebill.BillStatusCompleted {
		return nil, errors.NewValidationError("exit bill must be completed before creating entry bill")
	}

	// Check if entry bill already exists
	existingEntryBill, _ := s.warehouseBillRepo.FindByRelatedBillID(req.ExitBillID)
	if existingEntryBill != nil {
		return nil, errors.NewConflictError("entry bill already exists for this exit bill")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create entry bill items from exit bill items
	billItems := make([]warehousebill.WarehouseBillItem, len(exitBill.Items))
	totalAmount := 0.0
	for i, exitItem := range exitBill.Items {
		item := warehousebill.WarehouseBillItem{
			ProductVariantID: exitItem.ProductVariantID,
			ExpectedQuantity: exitItem.Quantity, // Copy expected quantity from exit bill
			Quantity:         exitItem.Quantity, // Keep for consistency
			UnitPrice:        exitItem.UnitPrice,
			DiscrepancyType:  warehousebill.DiscrepancyTypeNone,
		}
		item.CalculateTotal()
		billItems[i] = item
		totalAmount += item.TotalAmount
	}

	// Create entry bill
	bill := &warehousebill.WarehouseBill{
		CompanyID:          exitBill.CompanyID,
		FranchiseID:        franchiseID,
		BillType:           warehousebill.BillTypeEntry,
		RelatedBillID:      &exitBill.ID,
		Status:             warehousebill.BillStatusDraft,
		VerificationStatus: warehousebill.VerificationStatusPending,
		TotalAmount:        totalAmount,
		Notes:              req.Notes,
		CreatedByID:        userID,
		Items:              billItems,
	}

	if !bill.IsValid() {
		tx.Rollback()
		return nil, errors.NewValidationError("invalid bill data")
	}

	if err := s.warehouseBillRepo.Create(bill); err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to create entry bill", err)
	}

	tx.Commit()
	
	// Enrich response with product/variant details
	response := ToWarehouseBillResponse(bill)
	if err := s.enrichWarehouseBillResponse(response); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to enrich bill response: %v\n", err)
	}
	
	// Send email notification when bill is verified
	s.sendWarehouseBillEmail(bill, "entry")
	
	return response, nil
}

// VerifyEntryBill verifies received items and records discrepancies
func (s *Service) VerifyEntryBill(userID, franchiseID, billID uint, req *VerifyEntryBillRequest) (*WarehouseBillResponse, error) {
	// Check user authorization
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}

	franchiseRole, err := s.userRepo.FindUserRoleInFranchise(userID, franchiseID)
	if err != nil || franchiseRole == nil {
		parentRole, err := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
		if err != nil || parentRole == nil {
			return nil, errors.NewForbiddenError("you don't have access to this franchise")
		}
	}

	// Find entry bill
	bill, err := s.warehouseBillRepo.FindByID(billID)
	if err != nil {
		return nil, errors.NewNotFoundError("warehouse bill not found")
	}

	// Validate bill
	if !bill.IsEntryBill() {
		return nil, errors.NewValidationError("can only verify entry bills")
	}
	if bill.FranchiseID != franchiseID {
		return nil, errors.NewForbiddenError("bill does not belong to this franchise")
	}
	if !bill.CanBeVerified() {
		return nil, errors.NewValidationError("bill cannot be verified in current status")
	}

	// Get exit bill for comparison
	var exitBill *warehousebill.WarehouseBill
	if bill.RelatedBillID != nil {
		exitBill, err = s.warehouseBillRepo.FindByID(*bill.RelatedBillID)
		if err != nil {
			return nil, errors.NewNotFoundError("related exit bill not found")
		}
	}

	// Create map of received items by product variant ID
	receivedItemsMap := make(map[uint]int)
	for _, receivedItem := range req.Items {
		receivedItemsMap[receivedItem.ProductVariantID] = receivedItem.ReceivedQuantity
	}

	// Create map of expected items from exit bill
	expectedItemsMap := make(map[uint]int)
	if exitBill != nil {
		for _, exitItem := range exitBill.Items {
			expectedItemsMap[exitItem.ProductVariantID] = exitItem.Quantity
		}
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	hasDiscrepancies := false

	// Process each item in the entry bill
	for i := range bill.Items {
		item := &bill.Items[i]
		expectedQty := item.ExpectedQuantity
		receivedQty, received := receivedItemsMap[item.ProductVariantID]

		if !received || receivedQty == 0 {
			// Missing item
			if expectedQty > 0 {
				item.SetDiscrepancy(
					warehousebill.DiscrepancyTypeMissing,
					fmt.Sprintf("Expected %d, received 0", expectedQty),
				)
				item.ReceivedQuantity = nil
				hasDiscrepancies = true
			}
		} else if receivedQty != expectedQty {
			// Quantity mismatch
			item.ReceivedQuantity = &receivedQty
			item.SetDiscrepancy(
				warehousebill.DiscrepancyTypeQuantityMismatch,
				fmt.Sprintf("Expected %d, received %d", expectedQty, receivedQty),
			)
			item.CalculateTotal()
			hasDiscrepancies = true
		} else {
			// No discrepancy
			item.ReceivedQuantity = &receivedQty
			item.SetDiscrepancy(warehousebill.DiscrepancyTypeNone, "")
			item.CalculateTotal()
		}
	}

	// Check for extra items (received but not in exit bill)
	for variantID, receivedQty := range receivedItemsMap {
		if receivedQty > 0 {
			// Check if this variant is in the entry bill
			found := false
			for _, item := range bill.Items {
				if item.ProductVariantID == variantID {
					found = true
					break
				}
			}
			if !found {
				// Extra item - add to bill items
				extraItem := warehousebill.WarehouseBillItem{
					ProductVariantID: variantID,
					ExpectedQuantity: 0,
					ReceivedQuantity: &receivedQty,
					Quantity:         0,
					UnitPrice:        0, // Will need to get from product variant
					DiscrepancyType:  warehousebill.DiscrepancyTypeExtra,
				}
				// Get unit price from exit bill if available, or from product variant
				if exitBill != nil {
					for _, exitItem := range exitBill.Items {
						if exitItem.ProductVariantID == variantID {
							extraItem.UnitPrice = exitItem.UnitPrice
							break
						}
					}
				}
				if extraItem.UnitPrice == 0 {
				// Try to get from product variant
				variant, err := s.productRepo.FindProductVariantByID(variantID)
				if err == nil && variant.RetailPrice != nil {
					extraItem.UnitPrice = *variant.RetailPrice
				}
				}
				extraItem.SetDiscrepancy(
					warehousebill.DiscrepancyTypeExtra,
					fmt.Sprintf("Extra item: received %d (not expected)", receivedQty),
				)
				extraItem.CalculateTotal()
				bill.Items = append(bill.Items, extraItem)
				bill.TotalAmount += extraItem.TotalAmount
				hasDiscrepancies = true
			}
		}
	}

	// Recalculate total amount
	bill.TotalAmount = 0
	for _, item := range bill.Items {
		bill.TotalAmount += item.TotalAmount
	}

	// Update verification status
	if hasDiscrepancies {
		bill.MarkDiscrepanciesFound(userID)
	} else {
		bill.MarkVerified(userID)
	}

	// Save bill and items
	if err := s.warehouseBillRepo.Update(bill); err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update bill", err)
	}

	tx.Commit()
	
	// Enrich response with product/variant details
	response := ToWarehouseBillResponse(bill)
	if err := s.enrichWarehouseBillResponse(response); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to enrich bill response: %v\n", err)
	}
	
	return response, nil
}

// CompleteExitBill completes an exit bill and updates inventory
func (s *Service) CompleteExitBill(userID, companyID, billID uint) (*WarehouseBillResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, userDomain.RoleManager); err != nil {
		return nil, err
	}

	// Find bill
	bill, err := s.warehouseBillRepo.FindByID(billID)
	if err != nil {
		return nil, errors.NewNotFoundError("warehouse bill not found")
	}

	// Validate bill
	if !bill.IsExitBill() {
		return nil, errors.NewValidationError("can only complete exit bills")
	}
	if bill.CompanyID != companyID {
		return nil, errors.NewForbiddenError("bill does not belong to this company")
	}
	if !bill.CanBeCompleted() {
		return nil, errors.NewValidationError("bill cannot be completed in current status")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update inventory for each item
	for _, item := range bill.Items {
		// Get variant and product information for better error messages
		variant, err := s.productRepo.FindProductVariantByID(item.ProductVariantID)
		if err != nil {
			tx.Rollback()
			return nil, errors.NewNotFoundError(fmt.Sprintf("product variant %d not found", item.ProductVariantID))
		}

		product, err := s.productRepo.FindProductByID(variant.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, errors.NewNotFoundError(fmt.Sprintf("product not found for variant SKU: %s", variant.SKU))
		}

		// Find company inventory within transaction
		var companyInv inventory.Inventory
		if err := tx.Where("product_variant_id = ? AND company_id = ?", item.ProductVariantID, companyID).First(&companyInv).Error; err != nil {
			tx.Rollback()
			// Inventory not found means available quantity is 0
			return nil, errors.NewInternalError(fmt.Sprintf("inventory not found for product '%s' (SKU: %s), variant SKU: %s. Available quantity: 0", product.Name, product.SKU, variant.SKU), err)
		}

		// Stock is already reserved when bill was created as draft
		// On completion, we transfer the reserved stock (decrease reserved, decrease actual stock)
		// Verify stock is still reserved
		if companyInv.ReservedStock < item.Quantity {
			tx.Rollback()
			return nil, errors.NewValidationError(fmt.Sprintf("insufficient reserved stock for product '%s' (SKU: %s), variant SKU: %s. Reserved: %d, Required: %d", product.Name, product.SKU, variant.SKU, companyInv.ReservedStock, item.Quantity))
		}

		// Release reserved stock (it was reserved for this bill)
		previousReservedStock := companyInv.ReservedStock
		companyInv.ReleaseStock(item.Quantity)
		
		// Decrease actual stock (transfer to franchise)
		previousStock := companyInv.Stock
		if companyInv.Stock < item.Quantity {
			tx.Rollback()
			return nil, errors.NewValidationError(fmt.Sprintf("insufficient stock for product '%s' (SKU: %s), variant SKU: %s. Stock: %d, Required: %d", product.Name, product.SKU, variant.SKU, companyInv.Stock, item.Quantity))
		}
		companyInv.Stock -= item.Quantity

		// Update inventory
		if err := tx.Save(&companyInv).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to update inventory", err)
		}

		// Create inventory movement for transfer (completion)
		refType := "warehouse_bill"
		refID := fmt.Sprintf("%d", bill.ID)
		notes := fmt.Sprintf("Exit bill completed - stock transferred to franchise")
		movement := &inventory.InventoryMovement{
			InventoryID:   companyInv.ID,
			MovementType:  inventory.MovementTypeTransfer,
			Quantity:      item.Quantity,
			PreviousStock: previousStock,
			NewStock:      companyInv.Stock,
			ReferenceType: &refType,
			ReferenceID:   &refID,
			Notes:         &notes,
			CreatedByID:   userID,
		}
		if err := tx.Create(movement).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to create inventory movement", err)
		}

		// Create movement to track release of reserved stock
		reserveReleaseNotes := fmt.Sprintf("Reserved stock released on exit bill completion")
		reserveReleaseMovement := &inventory.InventoryMovement{
			InventoryID:   companyInv.ID,
			MovementType:  inventory.MovementTypeRelease,
			Quantity:      item.Quantity,
			PreviousStock: previousReservedStock,
			NewStock:      companyInv.ReservedStock,
			ReferenceType: &refType,
			ReferenceID:   &refID,
			Notes:         &reserveReleaseNotes,
			CreatedByID:   userID,
		}
		if err := tx.Create(reserveReleaseMovement).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to create inventory movement", err)
		}
	}

	// Mark bill as completed
	bill.Complete()
	if err := s.warehouseBillRepo.Update(bill); err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update bill", err)
	}

	tx.Commit()
	
	// Enrich response with product/variant details
	response := ToWarehouseBillResponse(bill)
	if err := s.enrichWarehouseBillResponse(response); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to enrich bill response: %v\n", err)
	}
	
	// Send email notification
	s.sendWarehouseBillEmail(bill, "exit")
	
	return response, nil
}

// CompleteEntryBill completes an entry bill and updates inventory based on actual received quantities
func (s *Service) CompleteEntryBill(userID, franchiseID, billID uint) (*WarehouseBillResponse, error) {
	// Check user authorization
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}

	franchiseRole, err := s.userRepo.FindUserRoleInFranchise(userID, franchiseID)
	if err != nil || franchiseRole == nil {
		parentRole, err := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
		if err != nil || parentRole == nil {
			return nil, errors.NewForbiddenError("you don't have access to this franchise")
		}
	}

	// Find bill
	bill, err := s.warehouseBillRepo.FindByID(billID)
	if err != nil {
		return nil, errors.NewNotFoundError("warehouse bill not found")
	}

	// Validate bill
	if !bill.IsEntryBill() {
		return nil, errors.NewValidationError("can only complete entry bills")
	}
	if bill.FranchiseID != franchiseID {
		return nil, errors.NewForbiddenError("bill does not belong to this franchise")
	}
	if bill.Status != warehousebill.BillStatusVerified {
		return nil, errors.NewValidationError("bill must be verified before completion")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update inventory for each item based on actual received quantities
	for _, item := range bill.Items {
		receivedQty := 0
		if item.ReceivedQuantity != nil {
			receivedQty = *item.ReceivedQuantity
		}

		// Skip missing items (receivedQty = 0)
		if receivedQty == 0 && item.DiscrepancyType == warehousebill.DiscrepancyTypeMissing {
			// Release reservation from company inventory
			companyInv, err := s.inventoryRepo.FindByVariantAndCompany(item.ProductVariantID, bill.CompanyID)
			if err == nil {
				previousStock := companyInv.Stock
				companyInv.ReleaseStock(item.ExpectedQuantity)
				if err := s.inventoryRepo.Update(companyInv); err == nil {
					// Create inventory movement for released stock
					refType := "warehouse_bill"
					refID := fmt.Sprintf("%d", bill.ID)
					movement := &inventory.InventoryMovement{
						InventoryID:   companyInv.ID,
						MovementType:  inventory.MovementTypeRelease,
						Quantity:      item.ExpectedQuantity,
						PreviousStock: previousStock,
						NewStock:      companyInv.Stock,
						ReferenceType: &refType,
						ReferenceID:   &refID,
						CreatedByID:   userID,
					}
					s.inventoryRepo.CreateMovement(movement)
				}
			}
			continue
		}

		// Process received items
		if receivedQty > 0 {
			// Release reservation from company inventory
			companyInv, err := s.inventoryRepo.FindByVariantAndCompany(item.ProductVariantID, bill.CompanyID)
			if err == nil {
				previousStock := companyInv.Stock
				companyInv.ReleaseStock(item.ExpectedQuantity)
				if err := s.inventoryRepo.Update(companyInv); err == nil {
					// Create inventory movement
					refType := "warehouse_bill"
					refID := fmt.Sprintf("%d", bill.ID)
					movement := &inventory.InventoryMovement{
						InventoryID:   companyInv.ID,
						MovementType:  inventory.MovementTypeTransfer,
						Quantity:      receivedQty,
						PreviousStock: previousStock,
						NewStock:      companyInv.Stock,
						ReferenceType: &refType,
						ReferenceID:   &refID,
						CreatedByID:   userID,
					}
					s.inventoryRepo.CreateMovement(movement)
				}
			}

			// Add to franchise inventory
			franchiseInv, err := s.inventoryRepo.FindByVariantAndFranchise(item.ProductVariantID, franchiseID)
			if err != nil {
				// Create new inventory if it doesn't exist
				franchiseInv = &inventory.Inventory{
					ProductVariantID: item.ProductVariantID,
					FranchiseID:      &franchiseID,
					Stock:            receivedQty,
					ReservedStock:    0,
					IsActive:         true,
				}
				if err := s.inventoryRepo.Create(franchiseInv); err != nil {
					tx.Rollback()
					return nil, errors.NewInternalError("failed to create franchise inventory", err)
				}
			} else {
				previousStock := franchiseInv.Stock
				franchiseInv.AddStock(receivedQty)
				if err := s.inventoryRepo.Update(franchiseInv); err != nil {
					tx.Rollback()
					return nil, errors.NewInternalError("failed to update franchise inventory", err)
				}

				// Create inventory movement
				refType := "warehouse_bill"
				refID := fmt.Sprintf("%d", bill.ID)
				movement := &inventory.InventoryMovement{
					InventoryID:   franchiseInv.ID,
					MovementType:  inventory.MovementTypeTransfer,
					Quantity:      receivedQty,
					PreviousStock: previousStock,
					NewStock:      franchiseInv.Stock,
					ReferenceType: &refType,
					ReferenceID:   &refID,
					CreatedByID:   userID,
				}
				if err := s.inventoryRepo.CreateMovement(movement); err != nil {
					tx.Rollback()
					return nil, errors.NewInternalError("failed to create inventory movement", err)
				}
			}
		}
	}

	// Mark bill as completed
	bill.Complete()
	if err := s.warehouseBillRepo.Update(bill); err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update bill", err)
	}

	tx.Commit()
	
	// Enrich response with product/variant details
	response := ToWarehouseBillResponse(bill)
	if err := s.enrichWarehouseBillResponse(response); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to enrich bill response: %v\n", err)
	}
	
	// Send email notification
	s.sendWarehouseBillEmail(bill, "entry")
	
	return response, nil
}

// sendWarehouseBillEmail sends an email notification for warehouse bill completion
func (s *Service) sendWarehouseBillEmail(bill *warehousebill.WarehouseBill, billType string) {
	if s.emailService == nil {
		return
	}

	// Get company
	company, err := s.companyRepo.FindByID(bill.CompanyID)
	if err != nil {
		fmt.Printf("Failed to get company for email: %v\n", err)
		return
	}

	// Get company users (managers and admins) to send emails to
	companyUsers, err := s.userRepo.FindByCompanyID(bill.CompanyID)
	if err != nil {
		fmt.Printf("Failed to get company users for email: %v\n", err)
		return
	}

	// Filter to managers and admins
	var recipients []string
	for _, user := range companyUsers {
		role, err := s.userRepo.FindUserRoleInCompany(user.ID, bill.CompanyID)
		if err == nil && (role.Role == userDomain.RoleOwner || role.Role == userDomain.RoleAdmin || role.Role == userDomain.RoleManager) {
			recipients = append(recipients, user.Email)
		}
	}

	if len(recipients) == 0 {
		return // No recipients
	}

	// Prepare bill items for email
	billItems := make([]map[string]interface{}, 0, len(bill.Items))
	for _, item := range bill.Items {
		variant, _ := s.productRepo.FindProductVariantByID(item.ProductVariantID)
		product, _ := s.productRepo.FindProductByID(variant.ProductID)
		
		productName := ""
		if product != nil {
			productName = product.Name
		} else if variant != nil {
			productName = variant.Name
		}
		
		billItems = append(billItems, map[string]interface{}{
			"product_name": productName,
			"quantity":     item.Quantity,
		})
	}

	// Send email
	emailReq := &emailApp.SendWarehouseBillEmailRequest{
		CompanyID:   bill.CompanyID,
		To:          recipients,
		BillNumber:  bill.BillNumber,
		BillType:    billType,
		BillDate:    bill.CreatedAt.Format("2006-01-02 15:04:05"),
		BillItems:   billItems,
		TotalAmount: 0, // Can be calculated if needed
		CompanyName: company.Name,
	}

	if err := s.emailService.SendWarehouseBillEmail(emailReq); err != nil {
		fmt.Printf("Failed to send warehouse bill email: %v\n", err)
	}
}

// enrichWarehouseBillResponse enriches bill response with product/variant details
func (s *Service) enrichWarehouseBillResponse(response *WarehouseBillResponse) error {
	if len(response.Items) == 0 {
		return nil
	}

	for i := range response.Items {
		item := &response.Items[i]
		
		// Fetch variant details
		variant, err := s.productRepo.FindProductVariantByID(item.ProductVariantID)
		if err != nil {
			// If variant not found, skip enrichment for this item
			continue
		}

		// Fetch product details
		product, err := s.productRepo.FindProductByID(variant.ProductID)
		if err != nil {
			// If product not found, still populate variant info
			variantName := variant.Name
			variantSKU := variant.SKU
			item.VariantName = &variantName
			item.VariantSKU = &variantSKU
			continue
		}

		// Populate all fields
		productName := product.Name
		variantName := variant.Name
		variantSKU := variant.SKU
		item.ProductName = &productName
		item.VariantName = &variantName
		item.VariantSKU = &variantSKU
	}

	return nil
}

// ListWarehouseBills lists warehouse bills with pagination and filters
func (s *Service) ListWarehouseBills(userID, companyID uint, page, limit int, filters *PaginationRequest) (*PaginatedResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, userDomain.RoleEmployee); err != nil {
		return nil, err
	}

	// Convert DTO filters to domain filters
	var billFilters *warehousebill.BillFilters
	if filters != nil {
		billFilters = &warehousebill.BillFilters{}
		if filters.FranchiseID != nil {
			billFilters.FranchiseID = filters.FranchiseID
		}
		if filters.Status != nil {
			status := warehousebill.BillStatus(*filters.Status)
			if status.IsValid() {
				billFilters.Status = &status
			}
		}
		if filters.BillType != nil {
			billType := warehousebill.BillType(*filters.BillType)
			if billType.IsValid() {
				billFilters.BillType = &billType
			}
		}
		if filters.DateFrom != nil {
			dateFrom, err := time.Parse("2006-01-02", *filters.DateFrom)
			if err == nil {
				billFilters.DateFrom = &dateFrom
			}
		}
		if filters.DateTo != nil {
			dateTo, err := time.Parse("2006-01-02", *filters.DateTo)
			if err == nil {
				billFilters.DateTo = &dateTo
			}
		}
	}

	bills, total, err := s.warehouseBillRepo.FindByCompanyIDWithFilters(companyID, page, limit, billFilters)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch warehouse bills", err)
	}

	responses := make([]*WarehouseBillResponse, len(bills))
	for i, bill := range bills {
		response := ToWarehouseBillResponse(bill)
		if err := s.enrichWarehouseBillResponse(response); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: failed to enrich bill response: %v\n", err)
		}
		responses[i] = response
	}

	return NewPaginatedResponse(responses, total, page, limit), nil
}

// ListFranchiseWarehouseBills lists warehouse bills for a franchise
func (s *Service) ListFranchiseWarehouseBills(userID, franchiseID uint, page, limit int) (*PaginatedResponse, error) {
	// Check user authorization
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}

	franchiseRole, err := s.userRepo.FindUserRoleInFranchise(userID, franchiseID)
	if err != nil || franchiseRole == nil {
		parentRole, err := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
		if err != nil || parentRole == nil {
			return nil, errors.NewForbiddenError("you don't have access to this franchise")
		}
	}

	bills, total, err := s.warehouseBillRepo.FindByFranchiseID(franchiseID, page, limit)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch warehouse bills", err)
	}

	responses := make([]*WarehouseBillResponse, len(bills))
	for i, bill := range bills {
		response := ToWarehouseBillResponse(bill)
		if err := s.enrichWarehouseBillResponse(response); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: failed to enrich bill response: %v\n", err)
		}
		responses[i] = response
	}

	return NewPaginatedResponse(responses, total, page, limit), nil
}

// GetWarehouseBill gets a single warehouse bill
func (s *Service) GetWarehouseBill(userID, companyID, billID uint) (*WarehouseBillResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, userDomain.RoleEmployee); err != nil {
		return nil, err
	}

	bill, err := s.warehouseBillRepo.FindByID(billID)
	if err != nil {
		return nil, errors.NewNotFoundError("warehouse bill not found")
	}

	if bill.CompanyID != companyID {
		return nil, errors.NewForbiddenError("bill does not belong to this company")
	}

	response := ToWarehouseBillResponse(bill)
	if err := s.enrichWarehouseBillResponse(response); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to enrich bill response: %v\n", err)
	}
	
	return response, nil
}

// GetFranchiseWarehouseBill gets a single warehouse bill for a franchise
func (s *Service) GetFranchiseWarehouseBill(userID, franchiseID, billID uint) (*WarehouseBillResponse, error) {
	// Check user authorization
	franchise, err := s.franchiseRepo.FindByID(franchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}

	franchiseRole, err := s.userRepo.FindUserRoleInFranchise(userID, franchiseID)
	if err != nil || franchiseRole == nil {
		parentRole, err := s.userRepo.FindUserRoleInCompany(userID, franchise.ParentCompanyID)
		if err != nil || parentRole == nil {
			return nil, errors.NewForbiddenError("you don't have access to this franchise")
		}
	}

	bill, err := s.warehouseBillRepo.FindByID(billID)
	if err != nil {
		return nil, errors.NewNotFoundError("warehouse bill not found")
	}

	if bill.FranchiseID != franchiseID {
		return nil, errors.NewForbiddenError("bill does not belong to this franchise")
	}

	response := ToWarehouseBillResponse(bill)
	if err := s.enrichWarehouseBillResponse(response); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to enrich bill response: %v\n", err)
	}
	
	return response, nil
}

// UpdateExitBillItems updates items in an exit bill (only for draft bills)
func (s *Service) UpdateExitBillItems(userID, companyID, billID uint, req *UpdateExitBillItemsRequest) (*WarehouseBillResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, userDomain.RoleManager); err != nil {
		return nil, err
	}

	// Find bill
	bill, err := s.warehouseBillRepo.FindByID(billID)
	if err != nil {
		return nil, errors.NewNotFoundError("warehouse bill not found")
	}

	// Validate bill
	if !bill.IsExitBill() {
		return nil, errors.NewValidationError("can only update exit bills")
	}
	if bill.CompanyID != companyID {
		return nil, errors.NewForbiddenError("bill does not belong to this company")
	}
	if bill.Status != warehousebill.BillStatusDraft {
		return nil, errors.NewValidationError("can only update draft bills")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create a map of existing items by ID for quick lookup
	existingItemsMap := make(map[uint]*warehousebill.WarehouseBillItem)
	for i := range bill.Items {
		item := &bill.Items[i]
		existingItemsMap[item.ID] = item
	}

	// Track which items are being kept
	newItemsSet := make(map[uint]bool)
	for _, itemReq := range req.Items {
		if itemReq.ID != nil {
			newItemsSet[*itemReq.ID] = true
		}
	}

	// Collect all variant IDs we need to fetch inventories for
	variantIDsSet := make(map[uint]bool)
	
	// Variants for deleted items
	for itemID, item := range existingItemsMap {
		if !newItemsSet[itemID] {
			variantIDsSet[item.ProductVariantID] = true
		}
	}
	
	// Variants for updated/new items
	for _, itemReq := range req.Items {
		variantIDsSet[itemReq.ProductVariantID] = true
		if itemReq.ID != nil {
			existingItem := existingItemsMap[*itemReq.ID]
			if existingItem != nil && existingItem.ProductVariantID != itemReq.ProductVariantID {
				// Variant changed, need old variant too
				variantIDsSet[existingItem.ProductVariantID] = true
			}
		}
	}

	// Batch fetch all inventories we need
	variantIDs := make([]uint, 0, len(variantIDsSet))
	for vid := range variantIDsSet {
		variantIDs = append(variantIDs, vid)
	}
	
	inventoryMap := make(map[uint]*inventory.Inventory)
	if len(variantIDs) > 0 {
		var inventories []*inventory.Inventory
		err = tx.Where("product_variant_id IN ? AND company_id = ?", variantIDs, companyID).Find(&inventories).Error
		if err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to fetch inventories", err)
		}
		for _, inv := range inventories {
			inventoryMap[inv.ProductVariantID] = inv
		}
	}

	// Batch fetch all product variants we need
	variantsMap := make(map[uint]*productDomain.ProductVariant)
	if len(variantIDs) > 0 {
		var variants []*productDomain.ProductVariant
		err = tx.Where("id IN ?", variantIDs).Find(&variants).Error
		if err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to fetch product variants", err)
		}
		for _, variant := range variants {
			variantsMap[variant.ID] = variant
		}
	}

	// Process deletions: Release reserved stock for items being removed
	for itemID, item := range existingItemsMap {
		if !newItemsSet[itemID] {
			// Item is being deleted, release reserved stock
			companyInv, exists := inventoryMap[item.ProductVariantID]
			if exists {
				previousReservedStock := companyInv.ReservedStock
				companyInv.ReleaseStock(item.Quantity)
				
				if err := tx.Save(companyInv).Error; err == nil {
					// Create inventory movement for released stock
					refType := "warehouse_bill"
					refID := fmt.Sprintf("%d", bill.ID)
					notes := fmt.Sprintf("Item removed from exit bill (Item ID: %d)", item.ID)
					if req.ChangeReason != "" {
						notes += " - " + req.ChangeReason
					}
					movement := &inventory.InventoryMovement{
						InventoryID:   companyInv.ID,
						MovementType:  inventory.MovementTypeRelease,
						Quantity:      item.Quantity,
						PreviousStock: previousReservedStock,
						NewStock:      companyInv.ReservedStock,
						ReferenceType: &refType,
						ReferenceID:   &refID,
						Notes:         &notes,
						CreatedByID:   userID,
					}
					if err := tx.Create(movement).Error; err != nil {
						tx.Rollback()
						return nil, errors.NewInternalError("failed to create inventory movement", err)
					}
				}
			}
		}
	}

	// Process updates and additions
	var updatedItems []warehousebill.WarehouseBillItem
	totalAmount := 0.0
	
	for _, itemReq := range req.Items {
		var item warehousebill.WarehouseBillItem
		
		if itemReq.ID != nil {
			// Update existing item
			existingItem, exists := existingItemsMap[*itemReq.ID]
			if !exists {
				tx.Rollback()
				return nil, errors.NewValidationError(fmt.Sprintf("item with id %d not found in bill", *itemReq.ID))
			}
			
			item = *existingItem
			oldQuantity := item.Quantity
			oldVariantID := item.ProductVariantID
			
			// Validate variant if changed
			if item.ProductVariantID != itemReq.ProductVariantID {
				if _, exists := variantsMap[itemReq.ProductVariantID]; !exists {
					tx.Rollback()
					return nil, errors.NewNotFoundError(fmt.Sprintf("product variant %d not found", itemReq.ProductVariantID))
				}
			}
			
			// Handle variant change: Release stock from old variant
			if oldVariantID != itemReq.ProductVariantID {
				oldCompanyInv, exists := inventoryMap[oldVariantID]
				if exists {
					previousReservedStock := oldCompanyInv.ReservedStock
					oldCompanyInv.ReleaseStock(oldQuantity)
					
					if err := tx.Save(oldCompanyInv).Error; err == nil {
						refType := "warehouse_bill"
						refID := fmt.Sprintf("%d", bill.ID)
						notes := fmt.Sprintf("Variant changed for item (Item ID: %d, Old Variant: %d)", item.ID, oldVariantID)
						if req.ChangeReason != "" {
							notes += " - " + req.ChangeReason
						}
						movement := &inventory.InventoryMovement{
							InventoryID:   oldCompanyInv.ID,
							MovementType:  inventory.MovementTypeRelease,
							Quantity:      oldQuantity,
							PreviousStock: previousReservedStock,
							NewStock:      oldCompanyInv.ReservedStock,
							ReferenceType: &refType,
							ReferenceID:   &refID,
							Notes:         &notes,
							CreatedByID:   userID,
						}
						if err := tx.Create(movement).Error; err != nil {
							tx.Rollback()
							return nil, errors.NewInternalError("failed to create inventory movement", err)
						}
					}
				}
			}
			
			// Update item fields
			item.ProductVariantID = itemReq.ProductVariantID
			item.Quantity = itemReq.Quantity
			item.UnitPrice = itemReq.UnitPrice
			item.CalculateTotal()
			
			// Handle inventory movement for quantity/variant changes
			variantChanged := oldVariantID != itemReq.ProductVariantID
			if variantChanged || oldQuantity != item.Quantity {
				companyInv, exists := inventoryMap[item.ProductVariantID]
				if !exists {
					tx.Rollback()
					return nil, errors.NewNotFoundError(fmt.Sprintf("inventory not found for variant %d", item.ProductVariantID))
				}
				
				previousReservedStock := companyInv.ReservedStock
				requiredQty := item.Quantity
				
				if variantChanged {
					// New variant - need to reserve full quantity
					if companyInv.GetAvailableStock() < requiredQty {
						tx.Rollback()
						return nil, errors.NewValidationError(fmt.Sprintf("insufficient stock for variant %d. Available: %d, Required: %d", item.ProductVariantID, companyInv.GetAvailableStock(), requiredQty))
					}
					
					if !companyInv.ReserveStock(requiredQty) {
						tx.Rollback()
						return nil, errors.NewValidationError(fmt.Sprintf("failed to reserve stock for variant %d", item.ProductVariantID))
					}
					
					if err := tx.Save(companyInv).Error; err != nil {
						tx.Rollback()
						return nil, errors.NewInternalError("failed to update inventory", err)
					}
					
					// Create inventory movement
					refType := "warehouse_bill"
					refID := fmt.Sprintf("%d", bill.ID)
					notes := fmt.Sprintf("Variant changed from %d to %d, quantity: %d (Item ID: %d)", oldVariantID, item.ProductVariantID, requiredQty, item.ID)
					if req.ChangeReason != "" {
						notes += " - " + req.ChangeReason
					}
					movement := &inventory.InventoryMovement{
						InventoryID:   companyInv.ID,
						MovementType:  inventory.MovementTypeReserve,
						Quantity:      requiredQty,
						PreviousStock: previousReservedStock,
						NewStock:      companyInv.ReservedStock,
						ReferenceType: &refType,
						ReferenceID:   &refID,
						Notes:         &notes,
						CreatedByID:   userID,
					}
					if err := tx.Create(movement).Error; err != nil {
						tx.Rollback()
						return nil, errors.NewInternalError("failed to create inventory movement", err)
					}
				} else if item.Quantity > oldQuantity {
					// Quantity increased - reserve additional stock
					additionalQty := item.Quantity - oldQuantity
					if companyInv.GetAvailableStock() < additionalQty {
						tx.Rollback()
						return nil, errors.NewValidationError(fmt.Sprintf("insufficient stock for variant %d. Available: %d, Required: %d", item.ProductVariantID, companyInv.GetAvailableStock(), additionalQty))
					}
					
					if !companyInv.ReserveStock(additionalQty) {
						tx.Rollback()
						return nil, errors.NewValidationError(fmt.Sprintf("failed to reserve stock for variant %d", item.ProductVariantID))
					}
					
					if err := tx.Save(companyInv).Error; err != nil {
						tx.Rollback()
						return nil, errors.NewInternalError("failed to update inventory", err)
					}
					
					// Create inventory movement
					refType := "warehouse_bill"
					refID := fmt.Sprintf("%d", bill.ID)
					notes := fmt.Sprintf("Quantity increased from %d to %d (Item ID: %d)", oldQuantity, item.Quantity, item.ID)
					if req.ChangeReason != "" {
						notes += " - " + req.ChangeReason
					}
					movement := &inventory.InventoryMovement{
						InventoryID:   companyInv.ID,
						MovementType:  inventory.MovementTypeReserve,
						Quantity:      additionalQty,
						PreviousStock: previousReservedStock,
						NewStock:      companyInv.ReservedStock,
						ReferenceType: &refType,
						ReferenceID:   &refID,
						Notes:         &notes,
						CreatedByID:   userID,
					}
					if err := tx.Create(movement).Error; err != nil {
						tx.Rollback()
						return nil, errors.NewInternalError("failed to create inventory movement", err)
					}
				} else if item.Quantity < oldQuantity {
					// Quantity decreased - release excess stock
					excessQty := oldQuantity - item.Quantity
					companyInv.ReleaseStock(excessQty)
					
					if err := tx.Save(companyInv).Error; err != nil {
						tx.Rollback()
						return nil, errors.NewInternalError("failed to update inventory", err)
					}
					
					// Create inventory movement
					refType := "warehouse_bill"
					refID := fmt.Sprintf("%d", bill.ID)
					notes := fmt.Sprintf("Quantity decreased from %d to %d (Item ID: %d)", oldQuantity, item.Quantity, item.ID)
					if req.ChangeReason != "" {
						notes += " - " + req.ChangeReason
					}
					movement := &inventory.InventoryMovement{
						InventoryID:   companyInv.ID,
						MovementType:  inventory.MovementTypeRelease,
						Quantity:      excessQty,
						PreviousStock: previousReservedStock,
						NewStock:      companyInv.ReservedStock,
						ReferenceType: &refType,
						ReferenceID:   &refID,
						Notes:         &notes,
						CreatedByID:   userID,
					}
					if err := tx.Create(movement).Error; err != nil {
						tx.Rollback()
						return nil, errors.NewInternalError("failed to create inventory movement", err)
					}
				}
			}
		} else {
			// New item - validate variant and inventory
			variant, exists := variantsMap[itemReq.ProductVariantID]
			if !exists {
				tx.Rollback()
				return nil, errors.NewNotFoundError(fmt.Sprintf("product variant %d not found", itemReq.ProductVariantID))
			}
			
			companyInv, exists := inventoryMap[itemReq.ProductVariantID]
			if !exists {
				tx.Rollback()
				return nil, errors.NewNotFoundError(fmt.Sprintf("inventory not found for variant %d", itemReq.ProductVariantID))
			}
			
			if companyInv.GetAvailableStock() < itemReq.Quantity {
				tx.Rollback()
				return nil, errors.NewValidationError(fmt.Sprintf("insufficient stock for variant %d. Available: %d, Required: %d", itemReq.ProductVariantID, companyInv.GetAvailableStock(), itemReq.Quantity))
			}
			
			// Reserve stock for new item
			previousReservedStock := companyInv.ReservedStock
			if !companyInv.ReserveStock(itemReq.Quantity) {
				tx.Rollback()
				return nil, errors.NewValidationError(fmt.Sprintf("failed to reserve stock for variant %d", itemReq.ProductVariantID))
			}
			
			if err := tx.Save(companyInv).Error; err != nil {
				tx.Rollback()
				return nil, errors.NewInternalError("failed to update inventory", err)
			}
			
			// Create inventory movement
			refType := "warehouse_bill"
			refID := fmt.Sprintf("%d", bill.ID)
			notes := fmt.Sprintf("New item added to exit bill (Variant ID: %d, SKU: %s)", variant.ID, variant.SKU)
			if req.ChangeReason != "" {
				notes += " - " + req.ChangeReason
			}
			movement := &inventory.InventoryMovement{
				InventoryID:   companyInv.ID,
				MovementType:  inventory.MovementTypeReserve,
				Quantity:      itemReq.Quantity,
				PreviousStock: previousReservedStock,
				NewStock:      companyInv.ReservedStock,
				ReferenceType: &refType,
				ReferenceID:   &refID,
				Notes:         &notes,
				CreatedByID:   userID,
			}
			if err := tx.Create(movement).Error; err != nil {
				tx.Rollback()
				return nil, errors.NewInternalError("failed to create inventory movement", err)
			}
			
			// Create new item
			item = warehousebill.WarehouseBillItem{
				WarehouseBillID:   bill.ID,
				ProductVariantID:  itemReq.ProductVariantID,
				Quantity:          itemReq.Quantity,
				UnitPrice:         itemReq.UnitPrice,
				DiscrepancyType:   warehousebill.DiscrepancyTypeNone,
			}
			item.CalculateTotal()
		}
		
		updatedItems = append(updatedItems, item)
		totalAmount += item.TotalAmount
	}

	// Delete removed items from database
	for itemID := range existingItemsMap {
		if !newItemsSet[itemID] {
			if err := tx.Delete(&warehousebill.WarehouseBillItem{}, "id = ? AND warehouse_bill_id = ?", itemID, bill.ID).Error; err != nil {
				tx.Rollback()
				return nil, errors.NewInternalError("failed to delete item", err)
			}
		}
	}

	// Update or create items in database
	for i := range updatedItems {
		item := &updatedItems[i]
		if item.ID > 0 {
			// Update existing item
			if err := tx.Save(item).Error; err != nil {
				tx.Rollback()
				return nil, errors.NewInternalError("failed to update item", err)
			}
		} else {
			// Create new item
			item.WarehouseBillID = bill.ID
			if err := tx.Create(item).Error; err != nil {
				tx.Rollback()
				return nil, errors.NewInternalError("failed to create item", err)
			}
		}
	}

	// Update bill total
	bill.TotalAmount = totalAmount
	
	// Update bill in transaction
	if err := tx.Save(bill).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update bill", err)
	}

	tx.Commit()

	// Reload bill with all items (outside transaction for response)
	bill, err = s.warehouseBillRepo.FindByID(billID)
	if err != nil {
		return nil, errors.NewInternalError("failed to reload bill", err)
	}

	// Enrich response with product/variant details
	response := ToWarehouseBillResponse(bill)
	if err := s.enrichWarehouseBillResponse(response); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to enrich bill response: %v\n", err)
	}

	return response, nil
}

// CancelWarehouseBill cancels a warehouse bill
func (s *Service) CancelWarehouseBill(userID, companyID, billID uint) error {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, userDomain.RoleManager); err != nil {
		return err
	}

	bill, err := s.warehouseBillRepo.FindByID(billID)
	if err != nil {
		return errors.NewNotFoundError("warehouse bill not found")
	}

	if bill.CompanyID != companyID {
		return errors.NewForbiddenError("bill does not belong to this company")
	}

	if bill.Status != warehousebill.BillStatusDraft {
		return errors.NewValidationError("can only cancel draft bills. Current status: " + string(bill.Status))
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Release reserved stock for all items (only for exit bills)
	if bill.IsExitBill() {
		for _, item := range bill.Items {
			companyInv, err := s.inventoryRepo.FindByVariantAndCompany(item.ProductVariantID, companyID)
			if err == nil {
				previousReservedStock := companyInv.ReservedStock
				companyInv.ReleaseStock(item.Quantity)
				
				if err := s.inventoryRepo.Update(companyInv); err == nil {
					// Create inventory movement for released stock
					refType := "warehouse_bill"
					refID := fmt.Sprintf("%d", bill.ID)
					notes := fmt.Sprintf("Bill cancelled - released reserved stock (Item ID: %d)", item.ID)
					movement := &inventory.InventoryMovement{
						InventoryID:   companyInv.ID,
						MovementType:  inventory.MovementTypeRelease,
						Quantity:      item.Quantity,
						PreviousStock: previousReservedStock,
						NewStock:      companyInv.ReservedStock,
						ReferenceType: &refType,
						ReferenceID:   &refID,
						Notes:         &notes,
						CreatedByID:   userID,
					}
					s.inventoryRepo.CreateMovement(movement)
				}
			}
		}
	}

	// Mark bill as cancelled
	bill.Cancel()
	if err := s.warehouseBillRepo.Update(bill); err != nil {
		tx.Rollback()
		return errors.NewInternalError("failed to cancel bill", err)
	}

	tx.Commit()
	return nil
}

// SearchProductsForExitBill searches for products/variants with franchise pricing
func (s *Service) SearchProductsForExitBill(userID, companyID uint, req *SearchProductsRequest) ([]*ProductVariantSearchResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, userDomain.RoleManager); err != nil {
		return nil, err
	}

	// Validate franchise belongs to company
	franchise, err := s.franchiseRepo.FindByID(req.FranchiseID)
	if err != nil {
		return nil, errors.NewNotFoundError("franchise not found")
	}
	if franchise.ParentCompanyID != companyID {
		return nil, errors.NewForbiddenError("franchise does not belong to this company")
	}

	// Set default limit
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}

	// Search variants
	variants, err := s.productRepo.SearchVariantsByCompany(companyID, req.Query, limit)
	if err != nil {
		return nil, errors.NewInternalError("failed to search products", err)
	}

	// Build response with pricing information
	results := make([]*ProductVariantSearchResponse, 0, len(variants))
	for _, variant := range variants {
		// Get product (should be preloaded)
		var product *productDomain.Product
		if variant.Product != nil {
			product = variant.Product
		} else {
			// Fetch product if not preloaded
			product, err = s.productRepo.FindProductByID(variant.ProductID)
			if err != nil {
				continue // Skip if product not found
			}
		}

		// Get base pricing
		baseRetailPrice := product.BaseRetailPrice
		baseWholesalePrice := product.BaseWholesalePrice

		// Get variant-specific pricing
		var variantRetailPrice, variantWholesalePrice *float64
		if !variant.UseParentPricing {
			if variant.RetailPrice != nil && *variant.RetailPrice > 0 {
				variantRetailPrice = variant.RetailPrice
			}
			if variant.WholesalePrice != nil && *variant.WholesalePrice > 0 {
				variantWholesalePrice = variant.WholesalePrice
			}
		}

		// Calculate effective pricing (variant or base)
		effectiveRetailPrice := variant.GetEffectiveRetailPrice(baseRetailPrice)
		effectiveWholesalePrice := variant.GetEffectiveWholesalePrice(baseWholesalePrice)

		// Get franchise pricing if available
		var franchiseRetailPrice, franchiseWholesalePrice *float64
		franchisePricing, err := s.franchiseRepo.FindPricing(req.FranchiseID, variant.ID)
		if err == nil && franchisePricing != nil {
			if franchisePricing.RetailPrice != nil && *franchisePricing.RetailPrice > 0 {
				franchiseRetailPrice = franchisePricing.RetailPrice
				effectiveRetailPrice = *franchisePricing.RetailPrice
			}
			if franchisePricing.WholesalePrice != nil && *franchisePricing.WholesalePrice > 0 {
				franchiseWholesalePrice = franchisePricing.WholesalePrice
				effectiveWholesalePrice = *franchisePricing.WholesalePrice
			}
		}

		results = append(results, &ProductVariantSearchResponse{
			VariantID:            variant.ID,
			VariantName:          variant.Name,
			VariantSKU:           variant.SKU,
			ProductID:            product.ID,
			ProductName:          product.Name,
			ProductSKU:           product.SKU,
			BaseRetailPrice:      baseRetailPrice,
			BaseWholesalePrice:   baseWholesalePrice,
			VariantRetailPrice:   variantRetailPrice,
			VariantWholesalePrice: variantWholesalePrice,
			FranchiseRetailPrice: franchiseRetailPrice,
			FranchiseWholesalePrice: franchiseWholesalePrice,
			EffectiveRetailPrice: effectiveRetailPrice,
			EffectiveWholesalePrice: effectiveWholesalePrice,
			UseParentPricing:     variant.UseParentPricing,
		})
	}

	return results, nil
}

// Helper functions
func (s *Service) checkUserCompanyAccess(userID, companyID uint, minRole userDomain.Role) error {
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

func (s *Service) hasSufficientRole(userRole, requiredRole userDomain.Role) bool {
	roleHierarchy := map[userDomain.Role]int{
		userDomain.RoleEmployee: 1,
		userDomain.RoleManager:  2,
		userDomain.RoleAdmin:    3,
		userDomain.RoleOwner:    4,
	}

	return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}

