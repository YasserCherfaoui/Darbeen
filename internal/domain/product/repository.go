package product

type Repository interface {
	// Product operations
	CreateProduct(product *Product) error
	FindProductByID(id uint) (*Product, error)
	FindProductByIDAndCompany(id, companyID uint) (*Product, error)
	FindProductBySKUAndCompany(sku string, companyID uint) (*Product, error)
	FindProductsByCompanyID(companyID uint, page, limit int) ([]*Product, int64, error)
	UpdateProduct(product *Product) error
	SoftDeleteProduct(id uint) error

	// Product variant operations
	CreateProductVariant(variant *ProductVariant) error
	FindProductVariantByID(id uint) (*ProductVariant, error)
	FindProductVariantByIDAndProduct(id, productID uint) (*ProductVariant, error)
	FindProductVariantBySKUAndProduct(sku string, productID uint) (*ProductVariant, error)
	FindProductVariantsByProductID(productID uint) ([]*ProductVariant, error)
	UpdateProductVariant(variant *ProductVariant) error
	SoftDeleteProductVariant(id uint) error
	SearchVariantsByCompany(companyID uint, query string, limit int) ([]*ProductVariant, error)

	// Stock operations
	UpdateVariantStock(variantID uint, newStock int) error
	AddVariantStock(variantID uint, amount int) error
	RemoveVariantStock(variantID uint, amount int) error
}
