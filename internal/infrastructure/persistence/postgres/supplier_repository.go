package postgres

import (
	"fmt"

	"github.com/YasserCherfaoui/darween/internal/domain/product"
	"github.com/YasserCherfaoui/darween/internal/domain/supplier"
	"gorm.io/gorm"
)

type supplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) supplier.Repository {
	return &supplierRepository{db: db}
}

// Supplier operations
func (r *supplierRepository) CreateSupplier(s *supplier.Supplier) error {
	return r.db.Create(s).Error
}

func (r *supplierRepository) FindSupplierByID(id uint) (*supplier.Supplier, error) {
	var s supplier.Supplier
	err := r.db.Where("id = ?", id).First(&s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, err
	}
	return &s, nil
}

func (r *supplierRepository) FindSupplierByIDAndCompany(id, companyID uint) (*supplier.Supplier, error) {
	var s supplier.Supplier
	err := r.db.Where("id = ? AND company_id = ?", id, companyID).First(&s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("supplier not found")
		}
		return nil, err
	}
	return &s, nil
}

func (r *supplierRepository) FindSuppliersByCompanyID(companyID uint, page, limit int) ([]*supplier.Supplier, int64, error) {
	var suppliers []*supplier.Supplier
	var total int64

	// Count total
	err := r.db.Model(&supplier.Supplier{}).Where("company_id = ? AND is_active = ?", companyID, true).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch suppliers
	err = r.db.Where("company_id = ? AND is_active = ?", companyID, true).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&suppliers).Error

	return suppliers, total, err
}

func (r *supplierRepository) UpdateSupplier(s *supplier.Supplier) error {
	return r.db.Save(s).Error
}

func (r *supplierRepository) SoftDeleteSupplier(id uint) error {
	return r.db.Model(&supplier.Supplier{}).Where("id = ?", id).Update("is_active", false).Error
}

// Product-Supplier relationship
func (r *supplierRepository) FindProductsBySupplier(supplierID uint) ([]*product.Product, error) {
	var products []*product.Product
	err := r.db.Where("supplier_id = ? AND is_active = ?", supplierID, true).Find(&products).Error
	return products, err
}

// SupplierBill operations
func (r *supplierRepository) CreateSupplierBill(bill *supplier.SupplierBill) error {
	return r.db.Create(bill).Error
}

func (r *supplierRepository) FindSupplierBillByID(id uint) (*supplier.SupplierBill, error) {
	var bill supplier.SupplierBill
	err := r.db.Where("id = ?", id).
		Preload("Items.ProductVariant.Product").
		Preload("Supplier").
		First(&bill).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("supplier bill not found")
		}
		return nil, err
	}
	return &bill, nil
}

func (r *supplierRepository) FindSupplierBillByIDAndCompany(id, companyID uint) (*supplier.SupplierBill, error) {
	var bill supplier.SupplierBill
	err := r.db.Where("id = ? AND company_id = ?", id, companyID).
		Preload("Items.ProductVariant.Product").
		Preload("Supplier").
		First(&bill).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("supplier bill not found")
		}
		return nil, err
	}
	return &bill, nil
}

func (r *supplierRepository) FindSupplierBillsBySupplier(supplierID, companyID uint, page, limit int) ([]*supplier.SupplierBill, int64, error) {
	var bills []*supplier.SupplierBill
	var total int64

	// Count total
	err := r.db.Model(&supplier.SupplierBill{}).
		Where("supplier_id = ? AND company_id = ?", supplierID, companyID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch bills
	err = r.db.Where("supplier_id = ? AND company_id = ?", supplierID, companyID).
		Preload("Items.ProductVariant.Product").
		Preload("Supplier").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&bills).Error

	return bills, total, err
}

func (r *supplierRepository) FindSupplierBillsByCompany(companyID uint, page, limit int) ([]*supplier.SupplierBill, int64, error) {
	var bills []*supplier.SupplierBill
	var total int64

	// Count total
	err := r.db.Model(&supplier.SupplierBill{}).
		Where("company_id = ?", companyID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch bills
	err = r.db.Where("company_id = ?", companyID).
		Preload("Items.ProductVariant.Product").
		Preload("Supplier").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&bills).Error

	return bills, total, err
}

func (r *supplierRepository) FindUnpaidBillsBySupplier(supplierID, companyID uint) ([]*supplier.SupplierBill, error) {
	var bills []*supplier.SupplierBill
	err := r.db.Where("supplier_id = ? AND company_id = ? AND payment_status IN ?",
		supplierID, companyID, []supplier.PaymentStatus{supplier.PaymentStatusUnpaid, supplier.PaymentStatusPartiallyPaid}).
		Order("created_at ASC"). // Oldest first for FIFO distribution
		Find(&bills).Error
	return bills, err
}

func (r *supplierRepository) UpdateSupplierBill(bill *supplier.SupplierBill) error {
	return r.db.Save(bill).Error
}

func (r *supplierRepository) DeleteSupplierBill(id uint) error {
	return r.db.Delete(&supplier.SupplierBill{}, id).Error
}

// SupplierBillItem operations
func (r *supplierRepository) CreateSupplierBillItem(item *supplier.SupplierBillItem) error {
	return r.db.Create(item).Error
}

func (r *supplierRepository) FindSupplierBillItemsByBill(billID uint) ([]*supplier.SupplierBillItem, error) {
	var items []*supplier.SupplierBillItem
	err := r.db.Where("supplier_bill_id = ?", billID).Find(&items).Error
	return items, err
}

func (r *supplierRepository) FindSupplierBillItemByID(id uint) (*supplier.SupplierBillItem, error) {
	var item supplier.SupplierBillItem
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("supplier bill item not found")
		}
		return nil, err
	}
	return &item, nil
}

func (r *supplierRepository) UpdateSupplierBillItem(item *supplier.SupplierBillItem) error {
	return r.db.Save(item).Error
}

func (r *supplierRepository) DeleteSupplierBillItem(id uint) error {
	return r.db.Delete(&supplier.SupplierBillItem{}, id).Error
}

// SupplierPayment operations
func (r *supplierRepository) CreateSupplierPayment(payment *supplier.SupplierPayment) error {
	return r.db.Create(payment).Error
}

func (r *supplierRepository) FindSupplierPaymentByID(id uint) (*supplier.SupplierPayment, error) {
	var payment supplier.SupplierPayment
	err := r.db.Where("id = ?", id).Preload("Distributions").First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("supplier payment not found")
		}
		return nil, err
	}
	return &payment, nil
}

func (r *supplierRepository) FindSupplierPaymentsBySupplier(supplierID, companyID uint, page, limit int) ([]*supplier.SupplierPayment, int64, error) {
	var payments []*supplier.SupplierPayment
	var total int64

	// Count total
	err := r.db.Model(&supplier.SupplierPayment{}).
		Where("supplier_id = ? AND company_id = ?", supplierID, companyID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch payments
	err = r.db.Where("supplier_id = ? AND company_id = ?", supplierID, companyID).
		Preload("Distributions").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&payments).Error

	return payments, total, err
}

func (r *supplierRepository) UpdateSupplierPayment(payment *supplier.SupplierPayment) error {
	return r.db.Save(payment).Error
}

// SupplierPaymentDistribution operations
func (r *supplierRepository) CreateSupplierPaymentDistribution(distribution *supplier.SupplierPaymentDistribution) error {
	return r.db.Create(distribution).Error
}

func (r *supplierRepository) FindPaymentDistributionsByPayment(paymentID uint) ([]*supplier.SupplierPaymentDistribution, error) {
	var distributions []*supplier.SupplierPaymentDistribution
	err := r.db.Where("supplier_payment_id = ?", paymentID).Find(&distributions).Error
	return distributions, err
}

func (r *supplierRepository) FindPaymentDistributionsByBill(billID uint) ([]*supplier.SupplierPaymentDistribution, error) {
	var distributions []*supplier.SupplierPaymentDistribution
	err := r.db.Where("supplier_bill_id = ?", billID).Find(&distributions).Error
	return distributions, err
}

// Outstanding balance calculation
func (r *supplierRepository) CalculateSupplierOutstandingBalance(supplierID, companyID uint) (float64, error) {
	var total float64
	err := r.db.Model(&supplier.SupplierBill{}).
		Where("supplier_id = ? AND company_id = ? AND payment_status IN ?",
			supplierID, companyID, []supplier.PaymentStatus{supplier.PaymentStatusUnpaid, supplier.PaymentStatusPartiallyPaid}).
		Select("COALESCE(SUM(pending_amount), 0)").
		Scan(&total).Error
	return total, err
}

