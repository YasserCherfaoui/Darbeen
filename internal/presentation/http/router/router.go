package router

import (
	"github.com/YasserCherfaoui/darween/internal/infrastructure/security"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/handler"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler          *handler.AuthHandler
	userHandler          *handler.UserHandler
	companyHandler       *handler.CompanyHandler
	subscriptionHandler  *handler.SubscriptionHandler
	productHandler       *handler.ProductHandler
	supplierHandler      *handler.SupplierHandler
	inventoryHandler     *handler.InventoryHandler
	franchiseHandler     *handler.FranchiseHandler
	posHandler           *handler.POSHandler
	warehouseBillHandler *handler.WarehouseBillHandler
	smtpConfigHandler    *handler.SMTPConfigHandler
	emailHandler         *handler.EmailHandler
	jwtManager           *security.JWTManager
}

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	companyHandler *handler.CompanyHandler,
	subscriptionHandler *handler.SubscriptionHandler,
	productHandler *handler.ProductHandler,
	supplierHandler *handler.SupplierHandler,
	inventoryHandler *handler.InventoryHandler,
	franchiseHandler *handler.FranchiseHandler,
	posHandler *handler.POSHandler,
	warehouseBillHandler *handler.WarehouseBillHandler,
	smtpConfigHandler *handler.SMTPConfigHandler,
	emailHandler *handler.EmailHandler,
	jwtManager *security.JWTManager,
) *Router {
	return &Router{
		authHandler:          authHandler,
		userHandler:          userHandler,
		companyHandler:       companyHandler,
		subscriptionHandler:  subscriptionHandler,
		productHandler:       productHandler,
		supplierHandler:      supplierHandler,
		inventoryHandler:     inventoryHandler,
		franchiseHandler:     franchiseHandler,
		posHandler:           posHandler,
		warehouseBillHandler: warehouseBillHandler,
		smtpConfigHandler:    smtpConfigHandler,
		emailHandler:         emailHandler,
		jwtManager:           jwtManager,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	// Apply global middleware
	engine.Use(middleware.CORS())

	// API v1 group
	v1 := engine.Group("/api/v1")

	// Public routes (Auth)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/password-reset/request", r.authHandler.RequestPasswordReset)
		auth.POST("/password-reset/confirm", r.authHandler.ConfirmPasswordReset)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(r.jwtManager))

	// User routes
	users := protected.Group("/users")
	{
		users.GET("/me", r.userHandler.GetMe)
		users.PUT("/me", r.userHandler.UpdateMe)
		users.GET("", r.userHandler.ListUsers)
	}

	// Company routes
	companies := protected.Group("/companies")
	{
		companies.POST("", r.companyHandler.CreateCompany)
		companies.GET("", r.companyHandler.ListCompanies)
		companies.GET("/:companyId", r.companyHandler.GetCompany)
		companies.PUT("/:companyId", r.companyHandler.UpdateCompany)
		companies.POST("/:companyId/users", r.companyHandler.AddUserToCompany)
		companies.GET("/:companyId/users", r.companyHandler.ListCompanyUsers)
		companies.PUT("/:companyId/users/:userId/role", r.companyHandler.UpdateUserRole)
		companies.DELETE("/:companyId/users/:userId", r.companyHandler.RemoveUserFromCompany)

		// Subscription routes nested under company
		companies.GET("/:companyId/subscription", r.subscriptionHandler.GetSubscription)
		companies.PUT("/:companyId/subscription", r.subscriptionHandler.UpdateSubscription)

		// SMTP config routes nested under company
		companies.POST("/:companyId/smtp-configs", r.smtpConfigHandler.CreateSMTPConfig)
		companies.GET("/:companyId/smtp-configs", r.smtpConfigHandler.ListSMTPConfigs)
		companies.GET("/:companyId/smtp-configs/:configId", r.smtpConfigHandler.GetSMTPConfig)
		companies.PUT("/:companyId/smtp-configs/:configId", r.smtpConfigHandler.UpdateSMTPConfig)
		companies.DELETE("/:companyId/smtp-configs/:configId", r.smtpConfigHandler.DeleteSMTPConfig)
		companies.PUT("/:companyId/smtp-configs/:configId/default", r.smtpConfigHandler.SetDefaultSMTPConfig)

		// Email routes nested under company
		companies.POST("/:companyId/emails/send", r.emailHandler.SendEmail)

		// Product routes nested under company
		companies.POST("/:companyId/products", r.productHandler.CreateProduct)
		companies.GET("/:companyId/products", r.productHandler.ListProducts)
		companies.GET("/:companyId/products/:productId", r.productHandler.GetProduct)
		companies.PUT("/:companyId/products/:productId", r.productHandler.UpdateProduct)
		companies.DELETE("/:companyId/products/:productId", r.productHandler.DeleteProduct)

		// Product variant routes nested under products
		companies.POST("/:companyId/products/:productId/variants", r.productHandler.CreateProductVariant)
		companies.POST("/:companyId/products/:productId/variants/bulk", r.productHandler.BulkCreateProductVariants)
		companies.GET("/:companyId/products/:productId/variants", r.productHandler.ListProductVariants)
		companies.GET("/:companyId/products/:productId/variants/:variantId", r.productHandler.GetProductVariant)
		companies.PUT("/:companyId/products/:productId/variants/:variantId", r.productHandler.UpdateProductVariant)
		companies.DELETE("/:companyId/products/:productId/variants/:variantId", r.productHandler.DeleteProductVariant)

		// Label generation routes for products and variants
		companies.GET("/:companyId/products/:productId/label", r.productHandler.GenerateProductLabel)
		companies.GET("/:companyId/products/:productId/variants/:variantId/label", r.productHandler.GenerateVariantLabel)
		companies.POST("/:companyId/products/labels/bulk", r.productHandler.GenerateBulkLabels)

		// Supplier routes nested under company
		companies.POST("/:companyId/suppliers", r.supplierHandler.CreateSupplier)
		companies.GET("/:companyId/suppliers", r.supplierHandler.ListSuppliers)
		companies.GET("/:companyId/suppliers/:supplierId", r.supplierHandler.GetSupplier)
		companies.PUT("/:companyId/suppliers/:supplierId", r.supplierHandler.UpdateSupplier)
		companies.DELETE("/:companyId/suppliers/:supplierId", r.supplierHandler.DeleteSupplier)
		companies.GET("/:companyId/suppliers/:supplierId/products", r.supplierHandler.GetSupplierProducts)
		companies.GET("/:companyId/suppliers/:supplierId/outstanding", r.supplierHandler.GetSupplierOutstandingBalance)
		companies.POST("/:companyId/suppliers/:supplierId/payments", r.supplierHandler.RecordSupplierPayment)

		// Supplier bill routes nested under suppliers
		companies.POST("/:companyId/suppliers/:supplierId/bills", r.supplierHandler.CreateSupplierBill)
		companies.GET("/:companyId/suppliers/:supplierId/bills", r.supplierHandler.ListSupplierBills)
		companies.GET("/:companyId/suppliers/:supplierId/bills/:billId", r.supplierHandler.GetSupplierBill)
		companies.PUT("/:companyId/suppliers/:supplierId/bills/:billId", r.supplierHandler.UpdateSupplierBill)
		companies.DELETE("/:companyId/suppliers/:supplierId/bills/:billId", r.supplierHandler.DeleteSupplierBill)

		// Bill endpoint without supplier ID (for inventory movement references)
		companies.GET("/:companyId/bills/:billId", r.supplierHandler.GetSupplierBillByID)

		// Supplier bill item routes nested under bills
		companies.POST("/:companyId/suppliers/:supplierId/bills/:billId/items", r.supplierHandler.AddBillItem)
		companies.PUT("/:companyId/suppliers/:supplierId/bills/:billId/items/:itemId", r.supplierHandler.UpdateBillItem)
		companies.DELETE("/:companyId/suppliers/:supplierId/bills/:billId/items/:itemId", r.supplierHandler.RemoveBillItem)

		// Franchise routes
		companies.POST("/:companyId/franchises", r.franchiseHandler.CreateFranchise)
		companies.GET("/:companyId/franchises", r.franchiseHandler.ListFranchises)

		// Inventory routes
		companies.GET("/:companyId/inventory", r.inventoryHandler.GetCompanyInventory)
		companies.POST("/:companyId/inventory/initialize", r.inventoryHandler.InitializeCompanyInventory)

		// POS routes
		companies.POST("/:companyId/pos/customers", r.posHandler.CreateCustomer)
		companies.GET("/:companyId/pos/customers", r.posHandler.ListCustomers)
		companies.GET("/:companyId/pos/customers/:customerId", r.posHandler.GetCustomer)
		companies.PUT("/:companyId/pos/customers/:customerId", r.posHandler.UpdateCustomer)
		companies.DELETE("/:companyId/pos/customers/:customerId", r.posHandler.DeleteCustomer)

		companies.POST("/:companyId/pos/sales", r.posHandler.CreateSale)
		companies.GET("/:companyId/pos/sales", r.posHandler.ListSales)
		companies.GET("/:companyId/pos/sales/:saleId", r.posHandler.GetSale)
		companies.GET("/:companyId/pos/sales/:saleId/receipt", r.posHandler.GenerateReceipt)
		companies.POST("/:companyId/pos/sales/:saleId/payments", r.posHandler.AddPayment)
		companies.POST("/:companyId/pos/sales/:saleId/refund", r.posHandler.ProcessRefund)

		companies.GET("/:companyId/pos/refunds", r.posHandler.ListRefunds)

		companies.POST("/:companyId/pos/cash-drawer/open", r.posHandler.OpenCashDrawer)
		companies.GET("/:companyId/pos/cash-drawer/active", r.posHandler.GetActiveCashDrawer)
		companies.PUT("/:companyId/pos/cash-drawer/:drawerId/close", r.posHandler.CloseCashDrawer)
		companies.GET("/:companyId/pos/cash-drawer", r.posHandler.ListCashDrawers)

		companies.POST("/:companyId/pos/reports/sales", r.posHandler.GetSalesReport)

		// Warehouse bill routes (exit bills)
		companies.POST("/:companyId/warehouse-bills/exit", r.warehouseBillHandler.CreateExitBill)
		companies.GET("/:companyId/warehouse-bills/search", r.warehouseBillHandler.SearchProductsForExitBill)
		companies.GET("/:companyId/warehouse-bills", r.warehouseBillHandler.ListWarehouseBills)
		companies.GET("/:companyId/warehouse-bills/:billId", r.warehouseBillHandler.GetWarehouseBill)
		companies.PUT("/:companyId/warehouse-bills/:billId/items", r.warehouseBillHandler.UpdateExitBillItems)
		companies.PUT("/:companyId/warehouse-bills/:billId/complete", r.warehouseBillHandler.CompleteExitBill)
		companies.DELETE("/:companyId/warehouse-bills/:billId", r.warehouseBillHandler.CancelWarehouseBill)
	}

	// Franchise-specific routes
	franchises := protected.Group("/franchises")
	{
		franchises.GET("/:franchiseId", r.franchiseHandler.GetFranchise)
		franchises.PUT("/:franchiseId", r.franchiseHandler.UpdateFranchise)
		franchises.POST("/:franchiseId/inventory/initialize", r.franchiseHandler.InitializeFranchiseInventory)
		franchises.GET("/:franchiseId/inventory", r.inventoryHandler.GetFranchiseInventory)
		franchises.GET("/:franchiseId/pricing", r.franchiseHandler.GetFranchisePricing)
		franchises.POST("/:franchiseId/pricing", r.franchiseHandler.SetFranchisePricing)
		franchises.POST("/:franchiseId/pricing/bulk", r.franchiseHandler.BulkSetFranchisePricing)
		franchises.DELETE("/:franchiseId/pricing/:variantId", r.franchiseHandler.DeleteFranchisePricing)
		franchises.POST("/:franchiseId/users", r.franchiseHandler.AddUserToFranchise)
		franchises.GET("/:franchiseId/users", r.franchiseHandler.ListFranchiseUsers)
		franchises.PUT("/:franchiseId/users/:userId/role", r.franchiseHandler.UpdateUserRole)
		franchises.DELETE("/:franchiseId/users/:userId", r.franchiseHandler.RemoveUserFromFranchise)

		// Franchise POS routes
		franchises.GET("/:franchiseId/pos/sales", r.posHandler.ListFranchiseSales)
		franchises.GET("/:franchiseId/pos/refunds", r.posHandler.ListFranchiseRefunds)
		franchises.GET("/:franchiseId/pos/cash-drawer/active", r.posHandler.GetActiveFranchiseCashDrawer)
		franchises.GET("/:franchiseId/pos/cash-drawer", r.posHandler.ListFranchiseCashDrawers)
		franchises.POST("/:franchiseId/pos/reports/sales", r.posHandler.GetFranchiseSalesReport)

		// Warehouse bill routes (entry bills)
		franchises.POST("/:franchiseId/warehouse-bills/entry", r.warehouseBillHandler.CreateEntryBill)
		franchises.GET("/:franchiseId/warehouse-bills", r.warehouseBillHandler.ListFranchiseWarehouseBills)
		franchises.GET("/:franchiseId/warehouse-bills/:billId", r.warehouseBillHandler.GetFranchiseWarehouseBill)
		franchises.POST("/:franchiseId/warehouse-bills/:billId/verify", r.warehouseBillHandler.VerifyEntryBill)
		franchises.PUT("/:franchiseId/warehouse-bills/:billId/complete", r.warehouseBillHandler.CompleteEntryBill)
	}

	// Inventory routes
	inventory := protected.Group("/inventory")
	{
		inventory.POST("", r.inventoryHandler.CreateInventory)
		inventory.PUT("/:inventoryId/stock", r.inventoryHandler.UpdateInventoryStock)
		inventory.POST("/:inventoryId/stock/adjust", r.inventoryHandler.AdjustInventoryStock)
		inventory.POST("/:inventoryId/reserve", r.inventoryHandler.ReserveStock)
		inventory.POST("/:inventoryId/release", r.inventoryHandler.ReleaseStock)
		inventory.GET("/:inventoryId/movements", r.inventoryHandler.GetInventoryMovements)
	}

	// Health check
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
