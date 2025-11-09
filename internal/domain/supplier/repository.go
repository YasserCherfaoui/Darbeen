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
}

