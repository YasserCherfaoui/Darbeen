package pos

import "time"

// CustomerRepository defines the interface for customer data operations
type CustomerRepository interface {
	Create(customer *Customer) error
	Update(customer *Customer) error
	FindByID(id uint) (*Customer, error)
	FindByCompanyID(companyID uint, page, limit int) ([]*Customer, int64, error)
	FindByEmail(email string, companyID uint) (*Customer, error)
	Delete(id uint) error
}

// SaleRepository defines the interface for sale data operations
type SaleRepository interface {
	Create(sale *Sale) error
	Update(sale *Sale) error
	FindByID(id uint) (*Sale, error)
	FindByReceiptNumber(receiptNumber string) (*Sale, error)
	FindByCompanyID(companyID uint, page, limit int) ([]*Sale, int64, error)
	FindByFranchiseID(franchiseID uint, page, limit int) ([]*Sale, int64, error)
	FindByDateRange(companyID uint, franchiseID *uint, startDate, endDate time.Time, page, limit int) ([]*Sale, int64, error)
	FindByCustomerID(customerID uint, page, limit int) ([]*Sale, int64, error)
	GetSalesReport(companyID uint, franchiseID *uint, startDate, endDate time.Time) (*SalesReportData, error)
}

// SaleItemRepository defines the interface for sale item data operations
type SaleItemRepository interface {
	Create(item *SaleItem) error
	CreateBulk(items []SaleItem) error
	FindBySaleID(saleID uint) ([]SaleItem, error)
	Update(item *SaleItem) error
	Delete(id uint) error
}

// PaymentRepository defines the interface for payment data operations
type PaymentRepository interface {
	Create(payment *Payment) error
	Update(payment *Payment) error
	FindByID(id uint) (*Payment, error)
	FindBySaleID(saleID uint) ([]Payment, error)
	GetTotalPaidForSale(saleID uint) (float64, error)
}

// CashDrawerRepository defines the interface for cash drawer data operations
type CashDrawerRepository interface {
	Create(drawer *CashDrawer) error
	Update(drawer *CashDrawer) error
	FindByID(id uint) (*CashDrawer, error)
	FindActiveByCompanyID(companyID uint) (*CashDrawer, error)
	FindActiveByFranchiseID(franchiseID uint) (*CashDrawer, error)
	FindByCompanyID(companyID uint, page, limit int) ([]*CashDrawer, int64, error)
	FindByFranchiseID(franchiseID uint, page, limit int) ([]*CashDrawer, int64, error)
}

// CashDrawerTransactionRepository defines the interface for cash drawer transaction data operations
type CashDrawerTransactionRepository interface {
	Create(transaction *CashDrawerTransaction) error
	FindByCashDrawerID(drawerID uint) ([]CashDrawerTransaction, error)
	GetTotalByCashDrawerID(drawerID uint) (float64, error)
}

// RefundRepository defines the interface for refund data operations
type RefundRepository interface {
	Create(refund *Refund) error
	Update(refund *Refund) error
	FindByID(id uint) (*Refund, error)
	FindBySaleID(saleID uint) ([]*Refund, error)
	FindByCompanyID(companyID uint, page, limit int) ([]*Refund, int64, error)
	FindByFranchiseID(franchiseID uint, page, limit int) ([]*Refund, int64, error)
}

// SalesReportData represents aggregated sales data for reporting
type SalesReportData struct {
	TotalSales       int64
	TotalRevenue     float64
	TotalCash        float64
	TotalCard        float64
	TotalRefunded    float64
	AverageOrderValue float64
	SalesByDate      map[string]float64
}

