package postgres

import (
	"fmt"
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/warehousebill"
	"gorm.io/gorm"
)

type warehouseBillRepository struct {
	db *gorm.DB
}

func NewWarehouseBillRepository(db *gorm.DB) warehousebill.Repository {
	return &warehouseBillRepository{db: db}
}

// generateBillNumber generates a unique bill number
func (r *warehouseBillRepository) generateBillNumber(billType warehousebill.BillType, companyID, billID uint) string {
	timestamp := time.Now().Format("20060102")
	typePrefix := "EXIT"
	if billType == warehousebill.BillTypeEntry {
		typePrefix = "ENTRY"
	}
	return fmt.Sprintf("WB-%s-%d-%s-%d", typePrefix, companyID, timestamp, billID)
}

func (r *warehouseBillRepository) Create(bill *warehousebill.WarehouseBill) error {
	// First create the bill to get the ID
	if err := r.db.Create(bill).Error; err != nil {
		return err
	}

	// Generate and update bill number
	bill.BillNumber = r.generateBillNumber(bill.BillType, bill.CompanyID, bill.ID)
	return r.db.Save(bill).Error
}

func (r *warehouseBillRepository) FindByID(id uint) (*warehousebill.WarehouseBill, error) {
	var bill warehousebill.WarehouseBill
	err := r.db.Where("id = ?", id).
		Preload("Items").
		First(&bill).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("warehouse bill not found")
		}
		return nil, err
	}
	return &bill, nil
}

func (r *warehouseBillRepository) FindByBillNumber(billNumber string) (*warehousebill.WarehouseBill, error) {
	var bill warehousebill.WarehouseBill
	err := r.db.Where("bill_number = ?", billNumber).
		Preload("Items").
		First(&bill).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("warehouse bill not found")
		}
		return nil, err
	}
	return &bill, nil
}

func (r *warehouseBillRepository) FindByCompanyID(companyID uint, page, limit int) ([]*warehousebill.WarehouseBill, int64, error) {
	return r.FindByCompanyIDWithFilters(companyID, page, limit, nil)
}

func (r *warehouseBillRepository) FindByCompanyIDWithFilters(companyID uint, page, limit int, filters *warehousebill.BillFilters) ([]*warehousebill.WarehouseBill, int64, error) {
	var bills []*warehousebill.WarehouseBill
	var total int64

	// Build query
	query := r.db.Model(&warehousebill.WarehouseBill{}).Where("company_id = ?", companyID)
	countQuery := r.db.Model(&warehousebill.WarehouseBill{}).Where("company_id = ?", companyID)

	// Apply filters
	if filters != nil {
		if filters.FranchiseID != nil {
			query = query.Where("franchise_id = ?", *filters.FranchiseID)
			countQuery = countQuery.Where("franchise_id = ?", *filters.FranchiseID)
		}
		if filters.Status != nil {
			query = query.Where("status = ?", *filters.Status)
			countQuery = countQuery.Where("status = ?", *filters.Status)
		}
		if filters.BillType != nil {
			query = query.Where("bill_type = ?", *filters.BillType)
			countQuery = countQuery.Where("bill_type = ?", *filters.BillType)
		}
		if filters.DateFrom != nil {
			query = query.Where("created_at >= ?", *filters.DateFrom)
			countQuery = countQuery.Where("created_at >= ?", *filters.DateFrom)
		}
		if filters.DateTo != nil {
			// Add one day to include the entire day
			dateTo := filters.DateTo.AddDate(0, 0, 1)
			query = query.Where("created_at < ?", dateTo)
			countQuery = countQuery.Where("created_at < ?", dateTo)
		}
	}

	// Count total
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch bills
	err = query.
		Preload("Items").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&bills).Error

	return bills, total, err
}

func (r *warehouseBillRepository) FindByFranchiseID(franchiseID uint, page, limit int) ([]*warehousebill.WarehouseBill, int64, error) {
	var bills []*warehousebill.WarehouseBill
	var total int64

	// Count total
	err := r.db.Model(&warehousebill.WarehouseBill{}).
		Where("franchise_id = ?", franchiseID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch bills
	err = r.db.Where("franchise_id = ?", franchiseID).
		Preload("Items").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&bills).Error

	return bills, total, err
}

func (r *warehouseBillRepository) FindByRelatedBillID(exitBillID uint) (*warehousebill.WarehouseBill, error) {
	var bill warehousebill.WarehouseBill
	err := r.db.Where("related_bill_id = ?", exitBillID).
		Preload("Items").
		First(&bill).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("entry bill not found for exit bill %d", exitBillID)
		}
		return nil, err
	}
	return &bill, nil
}

func (r *warehouseBillRepository) Update(bill *warehousebill.WarehouseBill) error {
	return r.db.Save(bill).Error
}

func (r *warehouseBillRepository) Delete(id uint) error {
	// Soft delete by setting status to cancelled
	return r.db.Model(&warehousebill.WarehouseBill{}).
		Where("id = ?", id).
		Update("status", warehousebill.BillStatusCancelled).Error
}

