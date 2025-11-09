package postgres

import (
	"time"

	"github.com/YasserCherfaoui/darween/internal/domain/pos"
	"gorm.io/gorm"
)

// CustomerRepositoryImpl implements the CustomerRepository interface
type CustomerRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) pos.CustomerRepository {
	return &CustomerRepositoryImpl{db: db}
}

func (r *CustomerRepositoryImpl) Create(customer *pos.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepositoryImpl) Update(customer *pos.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepositoryImpl) FindByID(id uint) (*pos.Customer, error) {
	var customer pos.Customer
	err := r.db.First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepositoryImpl) FindByCompanyID(companyID uint, page, limit int) ([]*pos.Customer, int64, error) {
	var customers []*pos.Customer
	var total int64

	query := r.db.Model(&pos.Customer{}).Where("company_id = ?", companyID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&customers).Error
	return customers, total, err
}

func (r *CustomerRepositoryImpl) FindByEmail(email string, companyID uint) (*pos.Customer, error) {
	var customer pos.Customer
	err := r.db.Where("email = ? AND company_id = ?", email, companyID).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&pos.Customer{}, id).Error
}

// SaleRepositoryImpl implements the SaleRepository interface
type SaleRepositoryImpl struct {
	db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) pos.SaleRepository {
	return &SaleRepositoryImpl{db: db}
}

func (r *SaleRepositoryImpl) Create(sale *pos.Sale) error {
	return r.db.Create(sale).Error
}

func (r *SaleRepositoryImpl) Update(sale *pos.Sale) error {
	return r.db.Save(sale).Error
}

func (r *SaleRepositoryImpl) FindByID(id uint) (*pos.Sale, error) {
	var sale pos.Sale
	err := r.db.Preload("Items").Preload("Payments").Preload("Customer").First(&sale, id).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *SaleRepositoryImpl) FindByReceiptNumber(receiptNumber string) (*pos.Sale, error) {
	var sale pos.Sale
	err := r.db.Preload("Items").Preload("Payments").Preload("Customer").
		Where("receipt_number = ?", receiptNumber).First(&sale).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *SaleRepositoryImpl) FindByCompanyID(companyID uint, page, limit int) ([]*pos.Sale, int64, error) {
	var sales []*pos.Sale
	var total int64

	query := r.db.Model(&pos.Sale{}).Where("company_id = ?", companyID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Items").Preload("Payments").Preload("Customer").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&sales).Error
	return sales, total, err
}

func (r *SaleRepositoryImpl) FindByFranchiseID(franchiseID uint, page, limit int) ([]*pos.Sale, int64, error) {
	var sales []*pos.Sale
	var total int64

	query := r.db.Model(&pos.Sale{}).Where("franchise_id = ?", franchiseID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Items").Preload("Payments").Preload("Customer").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&sales).Error
	return sales, total, err
}

func (r *SaleRepositoryImpl) FindByDateRange(companyID uint, franchiseID *uint, startDate, endDate time.Time, page, limit int) ([]*pos.Sale, int64, error) {
	var sales []*pos.Sale
	var total int64

	query := r.db.Model(&pos.Sale{}).
		Where("company_id = ? AND created_at >= ? AND created_at <= ?", companyID, startDate, endDate)

	if franchiseID != nil {
		query = query.Where("franchise_id = ?", *franchiseID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Items").Preload("Payments").Preload("Customer").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&sales).Error
	return sales, total, err
}

func (r *SaleRepositoryImpl) FindByCustomerID(customerID uint, page, limit int) ([]*pos.Sale, int64, error) {
	var sales []*pos.Sale
	var total int64

	query := r.db.Model(&pos.Sale{}).Where("customer_id = ?", customerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Items").Preload("Payments").Preload("Customer").
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&sales).Error
	return sales, total, err
}

func (r *SaleRepositoryImpl) GetSalesReport(companyID uint, franchiseID *uint, startDate, endDate time.Time) (*pos.SalesReportData, error) {
	query := r.db.Model(&pos.Sale{}).
		Where("company_id = ? AND created_at >= ? AND created_at <= ? AND sale_status = ?",
			companyID, startDate, endDate, pos.SaleStatusCompleted)

	if franchiseID != nil {
		query = query.Where("franchise_id = ?", *franchiseID)
	}

	var totalSales int64
	var totalRevenue float64
	
	err := query.Count(&totalSales).Error
	if err != nil {
		return nil, err
	}

	err = query.Select("COALESCE(SUM(total_amount), 0)").Row().Scan(&totalRevenue)
	if err != nil {
		return nil, err
	}

	// Get cash and card totals
	var totalCash, totalCard float64
	paymentQuery := r.db.Table("payments").
		Joins("JOIN sales ON payments.sale_id = sales.id").
		Where("sales.company_id = ? AND sales.created_at >= ? AND sales.created_at <= ? AND sales.sale_status = ? AND payments.payment_status = ?",
			companyID, startDate, endDate, pos.SaleStatusCompleted, pos.PaymentTransactionStatusCompleted)

	if franchiseID != nil {
		paymentQuery = paymentQuery.Where("sales.franchise_id = ?", *franchiseID)
	}

	paymentQuery.Where("payments.payment_method = ?", pos.PaymentMethodCash).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalCash)

	paymentQuery.Where("payments.payment_method = ?", pos.PaymentMethodCard).
		Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalCard)

	// Get refunded amount
	var totalRefunded float64
	refundQuery := r.db.Model(&pos.Refund{}).
		Joins("JOIN sales ON refunds.original_sale_id = sales.id").
		Where("sales.company_id = ? AND refunds.created_at >= ? AND refunds.created_at <= ? AND refunds.refund_status = ?",
			companyID, startDate, endDate, pos.RefundStatusCompleted)

	if franchiseID != nil {
		refundQuery = refundQuery.Where("sales.franchise_id = ?", *franchiseID)
	}

	refundQuery.Select("COALESCE(SUM(refund_amount), 0)").Row().Scan(&totalRefunded)

	averageOrderValue := 0.0
	if totalSales > 0 {
		averageOrderValue = totalRevenue / float64(totalSales)
	}

	// Get sales by date
	salesByDate := make(map[string]float64)
	type DateTotal struct {
		Date  string
		Total float64
	}
	var dateTotals []DateTotal
	
	query.Select("DATE(created_at) as date, SUM(total_amount) as total").
		Group("DATE(created_at)").
		Order("date").
		Scan(&dateTotals)

	for _, dt := range dateTotals {
		salesByDate[dt.Date] = dt.Total
	}

	return &pos.SalesReportData{
		TotalSales:        totalSales,
		TotalRevenue:      totalRevenue,
		TotalCash:         totalCash,
		TotalCard:         totalCard,
		TotalRefunded:     totalRefunded,
		AverageOrderValue: averageOrderValue,
		SalesByDate:       salesByDate,
	}, nil
}

// SaleItemRepositoryImpl implements the SaleItemRepository interface
type SaleItemRepositoryImpl struct {
	db *gorm.DB
}

func NewSaleItemRepository(db *gorm.DB) pos.SaleItemRepository {
	return &SaleItemRepositoryImpl{db: db}
}

func (r *SaleItemRepositoryImpl) Create(item *pos.SaleItem) error {
	return r.db.Create(item).Error
}

func (r *SaleItemRepositoryImpl) CreateBulk(items []pos.SaleItem) error {
	return r.db.Create(&items).Error
}

func (r *SaleItemRepositoryImpl) FindBySaleID(saleID uint) ([]pos.SaleItem, error) {
	var items []pos.SaleItem
	err := r.db.Where("sale_id = ?", saleID).Find(&items).Error
	return items, err
}

func (r *SaleItemRepositoryImpl) Update(item *pos.SaleItem) error {
	return r.db.Save(item).Error
}

func (r *SaleItemRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&pos.SaleItem{}, id).Error
}

// PaymentRepositoryImpl implements the PaymentRepository interface
type PaymentRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) pos.PaymentRepository {
	return &PaymentRepositoryImpl{db: db}
}

func (r *PaymentRepositoryImpl) Create(payment *pos.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepositoryImpl) Update(payment *pos.Payment) error {
	return r.db.Save(payment).Error
}

func (r *PaymentRepositoryImpl) FindByID(id uint) (*pos.Payment, error) {
	var payment pos.Payment
	err := r.db.First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) FindBySaleID(saleID uint) ([]pos.Payment, error) {
	var payments []pos.Payment
	err := r.db.Where("sale_id = ?", saleID).Find(&payments).Error
	return payments, err
}

func (r *PaymentRepositoryImpl) GetTotalPaidForSale(saleID uint) (float64, error) {
	var total float64
	err := r.db.Model(&pos.Payment{}).
		Where("sale_id = ? AND payment_status = ?", saleID, pos.PaymentTransactionStatusCompleted).
		Select("COALESCE(SUM(amount), 0)").
		Row().Scan(&total)
	return total, err
}

// CashDrawerRepositoryImpl implements the CashDrawerRepository interface
type CashDrawerRepositoryImpl struct {
	db *gorm.DB
}

func NewCashDrawerRepository(db *gorm.DB) pos.CashDrawerRepository {
	return &CashDrawerRepositoryImpl{db: db}
}

func (r *CashDrawerRepositoryImpl) Create(drawer *pos.CashDrawer) error {
	return r.db.Create(drawer).Error
}

func (r *CashDrawerRepositoryImpl) Update(drawer *pos.CashDrawer) error {
	return r.db.Save(drawer).Error
}

func (r *CashDrawerRepositoryImpl) FindByID(id uint) (*pos.CashDrawer, error) {
	var drawer pos.CashDrawer
	err := r.db.Preload("Transactions").First(&drawer, id).Error
	if err != nil {
		return nil, err
	}
	return &drawer, nil
}

func (r *CashDrawerRepositoryImpl) FindActiveByCompanyID(companyID uint) (*pos.CashDrawer, error) {
	var drawer pos.CashDrawer
	err := r.db.Preload("Transactions").
		Where("company_id = ? AND status = ?", companyID, pos.CashDrawerStatusOpen).
		First(&drawer).Error
	if err != nil {
		return nil, err
	}
	return &drawer, nil
}

func (r *CashDrawerRepositoryImpl) FindActiveByFranchiseID(franchiseID uint) (*pos.CashDrawer, error) {
	var drawer pos.CashDrawer
	err := r.db.Preload("Transactions").
		Where("franchise_id = ? AND status = ?", franchiseID, pos.CashDrawerStatusOpen).
		First(&drawer).Error
	if err != nil {
		return nil, err
	}
	return &drawer, nil
}

func (r *CashDrawerRepositoryImpl) FindByCompanyID(companyID uint, page, limit int) ([]*pos.CashDrawer, int64, error) {
	var drawers []*pos.CashDrawer
	var total int64

	query := r.db.Model(&pos.CashDrawer{}).Where("company_id = ?", companyID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Transactions").Offset(offset).Limit(limit).Order("opened_at DESC").Find(&drawers).Error
	return drawers, total, err
}

func (r *CashDrawerRepositoryImpl) FindByFranchiseID(franchiseID uint, page, limit int) ([]*pos.CashDrawer, int64, error) {
	var drawers []*pos.CashDrawer
	var total int64

	query := r.db.Model(&pos.CashDrawer{}).Where("franchise_id = ?", franchiseID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Transactions").Offset(offset).Limit(limit).Order("opened_at DESC").Find(&drawers).Error
	return drawers, total, err
}

// CashDrawerTransactionRepositoryImpl implements the CashDrawerTransactionRepository interface
type CashDrawerTransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewCashDrawerTransactionRepository(db *gorm.DB) pos.CashDrawerTransactionRepository {
	return &CashDrawerTransactionRepositoryImpl{db: db}
}

func (r *CashDrawerTransactionRepositoryImpl) Create(transaction *pos.CashDrawerTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *CashDrawerTransactionRepositoryImpl) FindByCashDrawerID(drawerID uint) ([]pos.CashDrawerTransaction, error) {
	var transactions []pos.CashDrawerTransaction
	err := r.db.Where("cash_drawer_id = ?", drawerID).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *CashDrawerTransactionRepositoryImpl) GetTotalByCashDrawerID(drawerID uint) (float64, error) {
	var total float64
	err := r.db.Model(&pos.CashDrawerTransaction{}).
		Where("cash_drawer_id = ?", drawerID).
		Select("COALESCE(SUM(CASE WHEN transaction_type = 'sale' THEN amount WHEN transaction_type = 'refund' THEN -amount ELSE amount END), 0)").
		Row().Scan(&total)
	return total, err
}

// RefundRepositoryImpl implements the RefundRepository interface
type RefundRepositoryImpl struct {
	db *gorm.DB
}

func NewRefundRepository(db *gorm.DB) pos.RefundRepository {
	return &RefundRepositoryImpl{db: db}
}

func (r *RefundRepositoryImpl) Create(refund *pos.Refund) error {
	return r.db.Create(refund).Error
}

func (r *RefundRepositoryImpl) Update(refund *pos.Refund) error {
	return r.db.Save(refund).Error
}

func (r *RefundRepositoryImpl) FindByID(id uint) (*pos.Refund, error) {
	var refund pos.Refund
	err := r.db.Preload("OriginalSale").First(&refund, id).Error
	if err != nil {
		return nil, err
	}
	return &refund, nil
}

func (r *RefundRepositoryImpl) FindBySaleID(saleID uint) ([]*pos.Refund, error) {
	var refunds []*pos.Refund
	err := r.db.Where("original_sale_id = ?", saleID).Find(&refunds).Error
	return refunds, err
}

func (r *RefundRepositoryImpl) FindByCompanyID(companyID uint, page, limit int) ([]*pos.Refund, int64, error) {
	var refunds []*pos.Refund
	var total int64

	query := r.db.Model(&pos.Refund{}).
		Joins("JOIN sales ON refunds.original_sale_id = sales.id").
		Where("sales.company_id = ?", companyID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("OriginalSale").Offset(offset).Limit(limit).Order("refunds.created_at DESC").Find(&refunds).Error
	return refunds, total, err
}

func (r *RefundRepositoryImpl) FindByFranchiseID(franchiseID uint, page, limit int) ([]*pos.Refund, int64, error) {
	var refunds []*pos.Refund
	var total int64

	query := r.db.Model(&pos.Refund{}).
		Joins("JOIN sales ON refunds.original_sale_id = sales.id").
		Where("sales.franchise_id = ?", franchiseID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("OriginalSale").Offset(offset).Limit(limit).Order("refunds.created_at DESC").Find(&refunds).Error
	return refunds, total, err
}

