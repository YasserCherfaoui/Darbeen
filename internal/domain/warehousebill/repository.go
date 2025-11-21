package warehousebill

import "time"

// BillFilters represents filters for querying warehouse bills
type BillFilters struct {
	FranchiseID *uint
	Status      *BillStatus
	BillType    *BillType
	DateFrom    *time.Time
	DateTo      *time.Time
}

type Repository interface {
	// Create creates a new warehouse bill
	Create(bill *WarehouseBill) error

	// FindByID finds a warehouse bill by ID
	FindByID(id uint) (*WarehouseBill, error)

	// FindByBillNumber finds a warehouse bill by bill number
	FindByBillNumber(billNumber string) (*WarehouseBill, error)

	// FindByCompanyID finds all warehouse bills for a company with pagination
	FindByCompanyID(companyID uint, page, limit int) ([]*WarehouseBill, int64, error)
	
	// FindByCompanyIDWithFilters finds warehouse bills for a company with filters and pagination
	FindByCompanyIDWithFilters(companyID uint, page, limit int, filters *BillFilters) ([]*WarehouseBill, int64, error)

	// FindByFranchiseID finds all warehouse bills for a franchise with pagination
	FindByFranchiseID(franchiseID uint, page, limit int) ([]*WarehouseBill, int64, error)

	// FindByRelatedBillID finds entry bill by related exit bill ID
	FindByRelatedBillID(exitBillID uint) (*WarehouseBill, error)

	// Update updates a warehouse bill
	Update(bill *WarehouseBill) error

	// Delete deletes a warehouse bill (soft delete by setting status to cancelled)
	Delete(id uint) error
}

