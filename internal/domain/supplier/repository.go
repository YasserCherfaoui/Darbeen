package supplier

import "github.com/YasserCherfaoui/darween/internal/domain/product"

type Repository interface {
	// Supplier operations
	CreateSupplier(supplier *Supplier) error
	FindSupplierByID(id uint) (*Supplier, error)
	FindSupplierByIDAndCompany(id, companyID uint) (*Supplier, error)
	FindSuppliersByCompanyID(companyID uint, page, limit int) ([]*Supplier, int64, error)
	UpdateSupplier(supplier *Supplier) error
	SoftDeleteSupplier(id uint) error

	// Product-Supplier relationship
	FindProductsBySupplier(supplierID uint) ([]*product.Product, error)

	// SupplierBill operations
	CreateSupplierBill(bill *SupplierBill) error
	FindSupplierBillByID(id uint) (*SupplierBill, error)
	FindSupplierBillByIDAndCompany(id, companyID uint) (*SupplierBill, error)
	FindSupplierBillsBySupplier(supplierID, companyID uint, page, limit int) ([]*SupplierBill, int64, error)
	FindSupplierBillsByCompany(companyID uint, page, limit int) ([]*SupplierBill, int64, error)
	FindUnpaidBillsBySupplier(supplierID, companyID uint) ([]*SupplierBill, error)
	UpdateSupplierBill(bill *SupplierBill) error
	DeleteSupplierBill(id uint) error

	// SupplierBillItem operations
	CreateSupplierBillItem(item *SupplierBillItem) error
	FindSupplierBillItemsByBill(billID uint) ([]*SupplierBillItem, error)
	FindSupplierBillItemByID(id uint) (*SupplierBillItem, error)
	UpdateSupplierBillItem(item *SupplierBillItem) error
	DeleteSupplierBillItem(id uint) error

	// SupplierPayment operations
	CreateSupplierPayment(payment *SupplierPayment) error
	FindSupplierPaymentByID(id uint) (*SupplierPayment, error)
	FindSupplierPaymentsBySupplier(supplierID, companyID uint, page, limit int) ([]*SupplierPayment, int64, error)
	UpdateSupplierPayment(payment *SupplierPayment) error

	// SupplierPaymentDistribution operations
	CreateSupplierPaymentDistribution(distribution *SupplierPaymentDistribution) error
	FindPaymentDistributionsByPayment(paymentID uint) ([]*SupplierPaymentDistribution, error)
	FindPaymentDistributionsByBill(billID uint) ([]*SupplierPaymentDistribution, error)

	// Outstanding balance calculation
	CalculateSupplierOutstandingBalance(supplierID, companyID uint) (float64, error)
}

