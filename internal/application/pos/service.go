package pos

import (
	"fmt"
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/company"
	"github.com/YasserCherfaoui/darween/internal/domain/franchise"
	"github.com/YasserCherfaoui/darween/internal/domain/inventory"
	"github.com/YasserCherfaoui/darween/internal/domain/pos"
	"github.com/YasserCherfaoui/darween/internal/domain/product"
	"github.com/YasserCherfaoui/darween/internal/domain/user"
	"github.com/YasserCherfaoui/darween/internal/infrastructure/receipt"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"gorm.io/gorm"
)

type Service struct {
	customerRepo              pos.CustomerRepository
	saleRepo                  pos.SaleRepository
	saleItemRepo              pos.SaleItemRepository
	paymentRepo               pos.PaymentRepository
	cashDrawerRepo            pos.CashDrawerRepository
	cashDrawerTransactionRepo pos.CashDrawerTransactionRepository
	refundRepo                pos.RefundRepository
	userRepo                  user.Repository
	inventoryRepo             inventory.Repository
	inventoryMovementRepo     inventory.Repository
	productVariantRepo        product.Repository
	franchiseRepo             franchise.Repository
	db                        *gorm.DB
}

func NewService(
	customerRepo pos.CustomerRepository,
	saleRepo pos.SaleRepository,
	saleItemRepo pos.SaleItemRepository,
	paymentRepo pos.PaymentRepository,
	cashDrawerRepo pos.CashDrawerRepository,
	cashDrawerTransactionRepo pos.CashDrawerTransactionRepository,
	refundRepo pos.RefundRepository,
	userRepo user.Repository,
	inventoryRepo inventory.Repository,
	inventoryMovementRepo inventory.Repository,
	productVariantRepo product.Repository,
	franchiseRepo franchise.Repository,
	db *gorm.DB,
) *Service {
	return &Service{
		customerRepo:              customerRepo,
		saleRepo:                  saleRepo,
		saleItemRepo:              saleItemRepo,
		paymentRepo:               paymentRepo,
		cashDrawerRepo:            cashDrawerRepo,
		cashDrawerTransactionRepo: cashDrawerTransactionRepo,
		refundRepo:                refundRepo,
		userRepo:                  userRepo,
		inventoryRepo:             inventoryRepo,
		inventoryMovementRepo:     inventoryMovementRepo,
		productVariantRepo:        productVariantRepo,
		franchiseRepo:             franchiseRepo,
		db:                        db,
	}
}

// Customer operations

func (s *Service) CreateCustomer(userID, companyID uint, req *CreateCustomerRequest) (*CustomerResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// Check if email already exists
	if req.Email != "" {
		existingCustomer, _ := s.customerRepo.FindByEmail(req.Email, companyID)
		if existingCustomer != nil {
			return nil, errors.NewConflictError("customer with this email already exists")
		}
	}

	customer := req.ToCustomer(companyID)
	if !customer.IsValid() {
		return nil, errors.NewValidationError("invalid customer data")
	}

	if err := s.customerRepo.Create(customer); err != nil {
		return nil, errors.NewInternalError("failed to create customer", err)
	}

	return ToCustomerResponse(customer), nil
}

func (s *Service) UpdateCustomer(userID, companyID, customerID uint, req *UpdateCustomerRequest) (*CustomerResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil {
		return nil, errors.NewNotFoundError("customer not found")
	}

	if customer.CompanyID != companyID {
		return nil, errors.NewForbiddenError("access denied to this customer")
	}

	// Update fields
	if req.Name != nil {
		customer.Name = *req.Name
	}
	if req.Email != nil {
		customer.Email = *req.Email
	}
	if req.Phone != nil {
		customer.Phone = *req.Phone
	}
	if req.Address != nil {
		customer.Address = *req.Address
	}
	if req.IsActive != nil {
		customer.IsActive = *req.IsActive
	}

	if err := s.customerRepo.Update(customer); err != nil {
		return nil, errors.NewInternalError("failed to update customer", err)
	}

	return ToCustomerResponse(customer), nil
}

func (s *Service) GetCustomerByID(userID, companyID, customerID uint) (*CustomerResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil {
		return nil, errors.NewNotFoundError("customer not found")
	}

	if customer.CompanyID != companyID {
		return nil, errors.NewForbiddenError("access denied to this customer")
	}

	return ToCustomerResponse(customer), nil
}

func (s *Service) ListCustomers(userID, companyID uint, page, limit int) (*PaginatedResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	customers, total, err := s.customerRepo.FindByCompanyID(companyID, page, limit)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch customers", err)
	}

	customerResponses := make([]*CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = ToCustomerResponse(customer)
	}

	return NewPaginatedResponse(customerResponses, total, page, limit), nil
}

func (s *Service) DeleteCustomer(userID, companyID, customerID uint) error {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return err
	}

	customer, err := s.customerRepo.FindByID(customerID)
	if err != nil {
		return errors.NewNotFoundError("customer not found")
	}

	if customer.CompanyID != companyID {
		return errors.NewForbiddenError("access denied to this customer")
	}

	if err := s.customerRepo.Delete(customerID); err != nil {
		return errors.NewInternalError("failed to delete customer", err)
	}

	return nil
}

// Sale operations

func (s *Service) CreateSale(userID, companyID uint, req *CreateSaleRequest) (*SaleResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// If franchise is specified, verify access
	if req.FranchiseID != nil {
		if err := s.checkUserFranchiseAccess(userID, *req.FranchiseID, user.RoleEmployee); err != nil {
			return nil, err
		}
	}

	// Validate customer if specified
	if req.CustomerID != nil {
		customer, err := s.customerRepo.FindByID(*req.CustomerID)
		if err != nil {
			return nil, errors.NewNotFoundError("customer not found")
		}
		if customer.CompanyID != companyID {
			return nil, errors.NewForbiddenError("customer does not belong to this company")
		}
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create sale
	sale := &pos.Sale{
		CompanyID:      companyID,
		FranchiseID:    req.FranchiseID,
		CustomerID:     req.CustomerID,
		TaxAmount:      req.TaxAmount,
		DiscountAmount: req.DiscountAmount,
		PaymentStatus:  pos.PaymentStatusUnpaid,
		SaleStatus:     pos.SaleStatusDraft,
		Notes:          req.Notes,
		CreatedByID:    userID,
	}

	// Create sale items and validate inventory
	saleItems := make([]pos.SaleItem, len(req.Items))
	for i, itemReq := range req.Items {
		// Validate product variant
		variant, err := s.productVariantRepo.FindProductVariantByID(itemReq.ProductVariantID)
		if err != nil {
			tx.Rollback()
			return nil, errors.NewNotFoundError(fmt.Sprintf("product variant %d not found", itemReq.ProductVariantID))
		}

		// Check inventory availability
		var inv *inventory.Inventory
		if req.FranchiseID != nil {
			inv, err = s.inventoryRepo.FindByVariantAndFranchise(itemReq.ProductVariantID, *req.FranchiseID)
		} else {
			inv, err = s.inventoryRepo.FindByVariantAndCompany(itemReq.ProductVariantID, companyID)
		}

		if err != nil || !inv.CanFulfill(itemReq.Quantity) {
			tx.Rollback()
			return nil, errors.NewValidationError(fmt.Sprintf("insufficient inventory for product variant %d (SKU: %s)", itemReq.ProductVariantID, variant.SKU))
		}

		saleItem := pos.SaleItem{
			ProductVariantID: itemReq.ProductVariantID,
			Quantity:         itemReq.Quantity,
			UnitPrice:        itemReq.UnitPrice,
			DiscountAmount:   itemReq.DiscountAmount,
		}
		saleItem.CalculateTotals()
		saleItems[i] = saleItem
	}

	// Calculate sale totals (temporarily set items for calculation)
	sale.Items = saleItems
	sale.CalculateTotals()
	sale.ReceiptNumber = GenerateReceiptNumber(companyID, 0) // Will update after getting ID

	if !sale.IsValid() {
		tx.Rollback()
		return nil, errors.NewValidationError("invalid sale data")
	}

	// Create sale in database without items to avoid GORM auto-creating them
	// We'll create items separately after getting the sale ID
	saleWithoutItems := *sale
	saleWithoutItems.Items = nil // Clear items to prevent GORM from creating them
	if err := tx.Create(&saleWithoutItems).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to create sale", err)
	}

	// Update receipt number with actual sale ID
	saleWithoutItems.ReceiptNumber = GenerateReceiptNumber(companyID, saleWithoutItems.ID)
	if err := tx.Save(&saleWithoutItems).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update receipt number", err)
	}

	// Update sale item IDs and create sale items
	for i := range saleItems {
		saleItems[i].SaleID = saleWithoutItems.ID
	}

	// Create sale items
	if err := tx.Create(&saleItems).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to create sale items", err)
	}

	// Update sale reference for later use
	sale.ID = saleWithoutItems.ID
	sale.ReceiptNumber = saleWithoutItems.ReceiptNumber

	// Deduct inventory and create movements
	for _, item := range saleItems {
		var inv *inventory.Inventory
		var err error
		if req.FranchiseID != nil {
			inv, err = s.inventoryRepo.FindByVariantAndFranchise(item.ProductVariantID, *req.FranchiseID)
		} else {
			inv, err = s.inventoryRepo.FindByVariantAndCompany(item.ProductVariantID, companyID)
		}

		if err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to fetch inventory", err)
		}

		previousStock := inv.Stock
		if !inv.RemoveStock(item.Quantity) {
			tx.Rollback()
			return nil, errors.NewValidationError(fmt.Sprintf("insufficient stock for variant %d", item.ProductVariantID))
		}

		if err := tx.Save(inv).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to update inventory", err)
		}

		// Create inventory movement
		saleIDStr := fmt.Sprintf("%d", sale.ID)
		movement := &inventory.InventoryMovement{
			InventoryID:   inv.ID,
			MovementType:  inventory.MovementTypeSale,
			Quantity:      -item.Quantity,
			PreviousStock: previousStock,
			NewStock:      inv.Stock,
			ReferenceType: stringPtr("sale"),
			ReferenceID:   &saleIDStr,
			CreatedByID:   userID,
		}

		if err := tx.Create(movement).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to create inventory movement", err)
		}
	}

	// Mark sale as completed
	saleWithoutItems.Complete()
	if err := tx.Save(&saleWithoutItems).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to complete sale", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.NewInternalError("failed to commit transaction", err)
	}

	// Reload sale with items
	completedSale, err := s.saleRepo.FindByID(sale.ID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch created sale", err)
	}

	return ToSaleResponse(completedSale), nil
}

func (s *Service) GetSaleByID(userID, companyID, saleID uint) (*SaleResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	sale, err := s.saleRepo.FindByID(saleID)
	if err != nil {
		return nil, errors.NewNotFoundError("sale not found")
	}

	if sale.CompanyID != companyID {
		return nil, errors.NewForbiddenError("access denied to this sale")
	}

	// Fetch product and variant details for each item
	itemDetails := make(map[uint]struct {
		productName string
		variantName string
		variantSKU  string
	})

	for _, item := range sale.Items {
		if _, exists := itemDetails[item.ProductVariantID]; exists {
			continue // Already fetched
		}

		variant, err := s.productVariantRepo.FindProductVariantByID(item.ProductVariantID)
		if err == nil && variant != nil {
			var productName string
			product, err := s.productVariantRepo.FindProductByID(variant.ProductID)
			if err == nil && product != nil {
				productName = product.Name
			}

			itemDetails[item.ProductVariantID] = struct {
				productName string
				variantName string
				variantSKU  string
			}{
				productName: productName,
				variantName: variant.Name,
				variantSKU:  variant.SKU,
			}
		}
	}

	return ToSaleResponseWithDetails(sale, itemDetails), nil
}

func (s *Service) ListSales(userID, companyID uint, franchiseID *uint, page, limit int) (*PaginatedResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	var sales []*pos.Sale
	var total int64
	var err error

	if franchiseID != nil {
		// Verify franchise access
		if err := s.checkUserFranchiseAccess(userID, *franchiseID, user.RoleEmployee); err != nil {
			return nil, err
		}
		sales, total, err = s.saleRepo.FindByFranchiseID(*franchiseID, page, limit)
	} else {
		sales, total, err = s.saleRepo.FindByCompanyID(companyID, page, limit)
	}

	if err != nil {
		return nil, errors.NewInternalError("failed to fetch sales", err)
	}

	saleResponses := make([]*SaleResponse, len(sales))
	for i, sale := range sales {
		saleResponses[i] = ToSaleResponse(sale)
	}

	return NewPaginatedResponse(saleResponses, total, page, limit), nil
}

// Payment operations

func (s *Service) AddPaymentToSale(userID, companyID, saleID uint, req *AddPaymentRequest) (*PaymentResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// Get sale
	sale, err := s.saleRepo.FindByID(saleID)
	if err != nil {
		return nil, errors.NewNotFoundError("sale not found")
	}

	if sale.CompanyID != companyID {
		return nil, errors.NewForbiddenError("access denied to this sale")
	}

	if sale.SaleStatus != pos.SaleStatusCompleted {
		return nil, errors.NewValidationError("can only add payments to completed sales")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create payment
	payment := &pos.Payment{
		SaleID:        saleID,
		PaymentMethod: req.PaymentMethod,
		Amount:        req.Amount,
		PaymentStatus: pos.PaymentTransactionStatusCompleted,
		Reference:     req.Reference,
		Notes:         req.Notes,
	}

	if !payment.IsValid() {
		tx.Rollback()
		return nil, errors.NewValidationError("invalid payment data")
	}

	if err := tx.Create(payment).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to create payment", err)
	}

	// Update sale payment status
	totalPaid, err := s.paymentRepo.GetTotalPaidForSale(saleID)
	if err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to calculate total paid", err)
	}

	sale.UpdatePaymentStatus(totalPaid)
	if err := tx.Save(sale).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update sale", err)
	}

	// If payment is cash, add to active cash drawer
	if req.PaymentMethod == pos.PaymentMethodCash {
		var activeDrawer *pos.CashDrawer
		var drawerErr error

		if sale.FranchiseID != nil {
			activeDrawer, drawerErr = s.cashDrawerRepo.FindActiveByFranchiseID(*sale.FranchiseID)
		} else {
			activeDrawer, drawerErr = s.cashDrawerRepo.FindActiveByCompanyID(companyID)
		}

		if drawerErr == nil && activeDrawer != nil {
			drawerTx := &pos.CashDrawerTransaction{
				CashDrawerID:    activeDrawer.ID,
				TransactionType: pos.CashDrawerTransactionTypeSale,
				Amount:          req.Amount,
				SaleID:          &saleID,
			}
			if err := tx.Create(drawerTx).Error; err != nil {
				// Don't fail payment if drawer transaction fails
				tx.Rollback()
				return nil, errors.NewInternalError("failed to record cash drawer transaction", err)
			}
		}
	}

	// Update customer total purchases if customer is linked
	if sale.CustomerID != nil && sale.PaymentStatus == pos.PaymentStatusPaid {
		customer, err := s.customerRepo.FindByID(*sale.CustomerID)
		if err == nil {
			customer.AddPurchase(sale.TotalAmount)
			tx.Save(customer)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.NewInternalError("failed to commit transaction", err)
	}

	return ToPaymentResponse(payment), nil
}

// Refund operations

func (s *Service) ProcessRefund(userID, companyID, saleID uint, req *ProcessRefundRequest) (*RefundResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleManager); err != nil {
		return nil, err
	}

	// Get sale
	sale, err := s.saleRepo.FindByID(saleID)
	if err != nil {
		return nil, errors.NewNotFoundError("sale not found")
	}

	if sale.CompanyID != companyID {
		return nil, errors.NewForbiddenError("access denied to this sale")
	}

	if !sale.CanBeRefunded() {
		return nil, errors.NewValidationError("sale cannot be refunded")
	}

	if req.RefundAmount > sale.TotalAmount {
		return nil, errors.NewValidationError("refund amount cannot exceed sale total")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create refund
	refund := &pos.Refund{
		OriginalSaleID: saleID,
		RefundAmount:   req.RefundAmount,
		Reason:         req.Reason,
		RefundMethod:   req.RefundMethod,
		RefundStatus:   pos.RefundStatusCompleted,
		ProcessedByID:  userID,
	}

	if !refund.IsValid() {
		tx.Rollback()
		return nil, errors.NewValidationError("invalid refund data")
	}

	if err := tx.Create(refund).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to create refund", err)
	}

	// Restore inventory for all items in the sale
	saleItems, err := s.saleItemRepo.FindBySaleID(saleID)
	if err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to fetch sale items", err)
	}

	for _, item := range saleItems {
		var inv *inventory.Inventory
		var invErr error
		if sale.FranchiseID != nil {
			inv, invErr = s.inventoryRepo.FindByVariantAndFranchise(item.ProductVariantID, *sale.FranchiseID)
		} else {
			inv, invErr = s.inventoryRepo.FindByVariantAndCompany(item.ProductVariantID, companyID)
		}

		if invErr != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to fetch inventory", invErr)
		}

		previousStock := inv.Stock
		inv.AddStock(item.Quantity)

		if err := tx.Save(inv).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to update inventory", err)
		}

		// Create inventory movement
		refundIDStr := fmt.Sprintf("%d", refund.ID)
		movement := &inventory.InventoryMovement{
			InventoryID:   inv.ID,
			MovementType:  inventory.MovementTypeReturn,
			Quantity:      item.Quantity,
			PreviousStock: previousStock,
			NewStock:      inv.Stock,
			ReferenceType: stringPtr("refund"),
			ReferenceID:   &refundIDStr,
			CreatedByID:   userID,
		}

		if err := tx.Create(movement).Error; err != nil {
			tx.Rollback()
			return nil, errors.NewInternalError("failed to create inventory movement", err)
		}
	}

	// Update sale status
	sale.SaleStatus = pos.SaleStatusRefunded
	sale.PaymentStatus = pos.PaymentStatusRefunded
	if err := tx.Save(sale).Error; err != nil {
		tx.Rollback()
		return nil, errors.NewInternalError("failed to update sale", err)
	}

	// If refund is cash, deduct from active cash drawer
	if req.RefundMethod == pos.PaymentMethodCash {
		var activeDrawer *pos.CashDrawer
		var drawerErr error

		if sale.FranchiseID != nil {
			activeDrawer, drawerErr = s.cashDrawerRepo.FindActiveByFranchiseID(*sale.FranchiseID)
		} else {
			activeDrawer, drawerErr = s.cashDrawerRepo.FindActiveByCompanyID(companyID)
		}

		if drawerErr == nil && activeDrawer != nil {
			drawerTx := &pos.CashDrawerTransaction{
				CashDrawerID:    activeDrawer.ID,
				TransactionType: pos.CashDrawerTransactionTypeRefund,
				Amount:          -req.RefundAmount,
				Notes:           fmt.Sprintf("Refund for sale #%s", sale.ReceiptNumber),
			}
			if err := tx.Create(drawerTx).Error; err != nil {
				tx.Rollback()
				return nil, errors.NewInternalError("failed to record cash drawer transaction", err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.NewInternalError("failed to commit transaction", err)
	}

	// Reload refund
	createdRefund, err := s.refundRepo.FindByID(refund.ID)
	if err != nil {
		return nil, errors.NewInternalError("failed to fetch created refund", err)
	}

	return ToRefundResponse(createdRefund), nil
}

func (s *Service) ListRefunds(userID, companyID uint, franchiseID *uint, page, limit int) (*PaginatedResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	var refunds []*pos.Refund
	var total int64
	var err error

	if franchiseID != nil {
		if err := s.checkUserFranchiseAccess(userID, *franchiseID, user.RoleEmployee); err != nil {
			return nil, err
		}
		refunds, total, err = s.refundRepo.FindByFranchiseID(*franchiseID, page, limit)
	} else {
		refunds, total, err = s.refundRepo.FindByCompanyID(companyID, page, limit)
	}

	if err != nil {
		return nil, errors.NewInternalError("failed to fetch refunds", err)
	}

	refundResponses := make([]*RefundResponse, len(refunds))
	for i, refund := range refunds {
		refundResponses[i] = ToRefundResponse(refund)
	}

	return NewPaginatedResponse(refundResponses, total, page, limit), nil
}

// Cash Drawer operations

func (s *Service) OpenCashDrawer(userID, companyID uint, req *OpenCashDrawerRequest) (*CashDrawerResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// If franchise is specified, verify access
	if req.FranchiseID != nil {
		if err := s.checkUserFranchiseAccess(userID, *req.FranchiseID, user.RoleEmployee); err != nil {
			return nil, err
		}
	}

	// Check if there's already an active drawer
	var activeDrawer *pos.CashDrawer
	var err error

	if req.FranchiseID != nil {
		activeDrawer, err = s.cashDrawerRepo.FindActiveByFranchiseID(*req.FranchiseID)
	} else {
		activeDrawer, err = s.cashDrawerRepo.FindActiveByCompanyID(companyID)
	}

	if err == nil && activeDrawer != nil {
		return nil, errors.NewConflictError("a cash drawer is already open")
	}

	// Create new cash drawer
	var companyIDPtr *uint
	if req.FranchiseID == nil {
		companyIDPtr = &companyID
	}

	drawer := &pos.CashDrawer{
		CompanyID:      companyIDPtr,
		FranchiseID:    req.FranchiseID,
		OpeningBalance: req.OpeningBalance,
		Status:         pos.CashDrawerStatusOpen,
		OpenedByID:     userID,
		OpenedAt:       time.Now(),
		Notes:          req.Notes,
	}

	if !drawer.IsValid() {
		return nil, errors.NewValidationError("invalid cash drawer data")
	}

	if err := s.cashDrawerRepo.Create(drawer); err != nil {
		return nil, errors.NewInternalError("failed to open cash drawer", err)
	}

	return ToCashDrawerResponse(drawer), nil
}

func (s *Service) CloseCashDrawer(userID, companyID uint, drawerID uint, req *CloseCashDrawerRequest) (*CashDrawerResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	drawer, err := s.cashDrawerRepo.FindByID(drawerID)
	if err != nil {
		return nil, errors.NewNotFoundError("cash drawer not found")
	}

	// Verify access
	if drawer.CompanyID != nil && *drawer.CompanyID != companyID {
		return nil, errors.NewForbiddenError("access denied to this cash drawer")
	}

	if drawer.FranchiseID != nil {
		if err := s.checkUserFranchiseAccess(userID, *drawer.FranchiseID, user.RoleEmployee); err != nil {
			return nil, err
		}
	}

	if !drawer.IsOpen() {
		return nil, errors.NewValidationError("cash drawer is already closed")
	}

	// Calculate expected balance
	totalTransactions, err := s.cashDrawerTransactionRepo.GetTotalByCashDrawerID(drawerID)
	if err != nil {
		return nil, errors.NewInternalError("failed to calculate transactions", err)
	}

	expectedBalance := drawer.OpeningBalance + totalTransactions

	// Close drawer
	drawer.Close(req.ClosingBalance, userID, expectedBalance)
	drawer.Notes = req.Notes

	if err := s.cashDrawerRepo.Update(drawer); err != nil {
		return nil, errors.NewInternalError("failed to close cash drawer", err)
	}

	return ToCashDrawerResponse(drawer), nil
}

func (s *Service) GetActiveCashDrawer(userID, companyID uint, franchiseID *uint) (*CashDrawerResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	var drawer *pos.CashDrawer
	var err error

	if franchiseID != nil {
		if err := s.checkUserFranchiseAccess(userID, *franchiseID, user.RoleEmployee); err != nil {
			return nil, err
		}
		drawer, err = s.cashDrawerRepo.FindActiveByFranchiseID(*franchiseID)
	} else {
		drawer, err = s.cashDrawerRepo.FindActiveByCompanyID(companyID)
	}

	if err != nil {
		return nil, errors.NewNotFoundError("no active cash drawer found")
	}

	return ToCashDrawerResponse(drawer), nil
}

func (s *Service) ListCashDrawers(userID, companyID uint, franchiseID *uint, page, limit int) (*PaginatedResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	var drawers []*pos.CashDrawer
	var total int64
	var err error

	if franchiseID != nil {
		if err := s.checkUserFranchiseAccess(userID, *franchiseID, user.RoleEmployee); err != nil {
			return nil, err
		}
		drawers, total, err = s.cashDrawerRepo.FindByFranchiseID(*franchiseID, page, limit)
	} else {
		drawers, total, err = s.cashDrawerRepo.FindByCompanyID(companyID, page, limit)
	}

	if err != nil {
		return nil, errors.NewInternalError("failed to fetch cash drawers", err)
	}

	drawerResponses := make([]*CashDrawerResponse, len(drawers))
	for i, drawer := range drawers {
		drawerResponses[i] = ToCashDrawerResponse(drawer)
	}

	return NewPaginatedResponse(drawerResponses, total, page, limit), nil
}

// Reports

func (s *Service) GetSalesReport(userID, companyID uint, req *SalesReportRequest) (*SalesReportResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	if req.FranchiseID != nil {
		if err := s.checkUserFranchiseAccess(userID, *req.FranchiseID, user.RoleEmployee); err != nil {
			return nil, err
		}
	}

	data, err := s.saleRepo.GetSalesReport(companyID, req.FranchiseID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, errors.NewInternalError("failed to generate sales report", err)
	}

	return ToSalesReportResponse(data), nil
}

// Receipt Generation

func (s *Service) GenerateReceiptPDF(userID, companyID, saleID uint) ([]byte, string, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, "", err
	}

	// Fetch sale with items, payments, and customer
	sale, err := s.saleRepo.FindByID(saleID)
	if err != nil {
		return nil, "", errors.NewNotFoundError("sale not found")
	}

	if sale.CompanyID != companyID {
		return nil, "", errors.NewForbiddenError("access denied to this sale")
	}

	// Fetch company info
	var company company.Company
	if err := s.db.First(&company, companyID).Error; err != nil {
		return nil, "", errors.NewNotFoundError("company not found")
	}

	// Fetch franchise info if applicable
	var franchiseName string
	if sale.FranchiseID != nil {
		var franchise franchise.Franchise
		if err := s.db.First(&franchise, *sale.FranchiseID).Error; err == nil {
			franchiseName = franchise.Name
		}
	}

	// Fetch product variants for items
	productVariantMap := make(map[uint]receipt.ProductVariantInfo)
	for _, item := range sale.Items {
		variant, err := s.productVariantRepo.FindProductVariantByID(item.ProductVariantID)
		if err == nil {
			// Get product name using ProductID from variant
			product, err := s.productVariantRepo.FindProductByID(variant.ProductID)
			if err == nil {
				itemName := fmt.Sprintf("%s - %s", product.Name, variant.Name)
				productVariantMap[item.ProductVariantID] = receipt.ProductVariantInfo{
					Name: itemName,
					SKU:  variant.SKU,
				}
			} else {
				// Fallback to variant name only
				productVariantMap[item.ProductVariantID] = receipt.ProductVariantInfo{
					Name: variant.Name,
					SKU:  variant.SKU,
				}
			}
		}
	}

	// Convert sale to receipt data
	receiptData := receipt.ConvertSaleToReceiptData(sale, company.Name, franchiseName, productVariantMap)

	// Generate PDF
	generator := receipt.NewReceiptGenerator()
	pdfBytes, err := generator.GenerateReceipt(receiptData)
	if err != nil {
		return nil, "", errors.NewInternalError("failed to generate receipt PDF", err)
	}

	// Generate filename
	filename := fmt.Sprintf("receipt_%s.pdf", sale.ReceiptNumber)

	return pdfBytes, filename, nil
}

// Helper methods

func (s *Service) checkUserCompanyAccess(userID, companyID uint, minimumRole user.Role) error {
	ucr, err := s.userRepo.FindUserRoleInCompany(userID, companyID)
	if err != nil {
		return errors.NewForbiddenError("access denied to this company")
	}

	if !ucr.Role.HasPermission(minimumRole) {
		return errors.NewForbiddenError("insufficient permissions")
	}

	return nil
}

func (s *Service) checkUserFranchiseAccess(userID, franchiseID uint, minimumRole user.Role) error {
	ufr, err := s.userRepo.FindUserRoleInFranchise(userID, franchiseID)
	if err != nil {
		return errors.NewForbiddenError("access denied to this franchise")
	}

	if !ufr.Role.HasPermission(minimumRole) {
		return errors.NewForbiddenError("insufficient permissions")
	}

	return nil
}

// SearchProductsForSale searches for products/variants with retail pricing for POS sales
func (s *Service) SearchProductsForSale(userID, companyID uint, req *SearchProductsRequest) ([]*ProductVariantSearchResponse, error) {
	// Check user authorization
	if err := s.checkUserCompanyAccess(userID, companyID, user.RoleEmployee); err != nil {
		return nil, err
	}

	// If franchise is specified, validate it belongs to company
	if req.FranchiseID != nil {
		franchise, err := s.franchiseRepo.FindByID(*req.FranchiseID)
		if err != nil {
			return nil, errors.NewNotFoundError("franchise not found")
		}
		if franchise.ParentCompanyID != companyID {
			return nil, errors.NewForbiddenError("franchise does not belong to this company")
		}
	}

	// Set default limit
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}

	// Search variants
	variants, err := s.productVariantRepo.SearchVariantsByCompany(companyID, req.Query, limit)
	if err != nil {
		return nil, errors.NewInternalError("failed to search products", err)
	}

	// Build response with retail pricing information
	results := make([]*ProductVariantSearchResponse, 0, len(variants))
	for _, variant := range variants {
		// Get product (should be preloaded)
		var product *product.Product
		if variant.Product != nil {
			product = variant.Product
		} else {
			// Fetch product if not preloaded
			product, err = s.productVariantRepo.FindProductByID(variant.ProductID)
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
		if req.FranchiseID != nil {
			franchisePricing, err := s.franchiseRepo.FindPricing(*req.FranchiseID, variant.ID)
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
		}

		results = append(results, &ProductVariantSearchResponse{
			VariantID:                variant.ID,
			VariantName:              variant.Name,
			VariantSKU:               variant.SKU,
			ProductID:                product.ID,
			ProductName:              product.Name,
			ProductSKU:               product.SKU,
			BaseRetailPrice:          baseRetailPrice,
			BaseWholesalePrice:       baseWholesalePrice,
			VariantRetailPrice:       variantRetailPrice,
			VariantWholesalePrice:    variantWholesalePrice,
			FranchiseRetailPrice:     franchiseRetailPrice,
			FranchiseWholesalePrice:  franchiseWholesalePrice,
			EffectiveRetailPrice:     effectiveRetailPrice,
			EffectiveWholesalePrice:   effectiveWholesalePrice,
			UseParentPricing:         variant.UseParentPricing,
		})
	}

	return results, nil
}

func stringPtr(s string) *string {
	return &s
}
