package router

import (
	"github.com/YasserCherfaoui/darween/internal/infrastructure/security"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/handler"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler         *handler.AuthHandler
	userHandler         *handler.UserHandler
	companyHandler      *handler.CompanyHandler
	subscriptionHandler *handler.SubscriptionHandler
	productHandler      *handler.ProductHandler
	inventoryHandler    *handler.InventoryHandler
	franchiseHandler    *handler.FranchiseHandler
	jwtManager          *security.JWTManager
}

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	companyHandler *handler.CompanyHandler,
	subscriptionHandler *handler.SubscriptionHandler,
	productHandler *handler.ProductHandler,
	inventoryHandler *handler.InventoryHandler,
	franchiseHandler *handler.FranchiseHandler,
	jwtManager *security.JWTManager,
) *Router {
	return &Router{
		authHandler:         authHandler,
		userHandler:         userHandler,
		companyHandler:      companyHandler,
		subscriptionHandler: subscriptionHandler,
		productHandler:      productHandler,
		inventoryHandler:    inventoryHandler,
		franchiseHandler:    franchiseHandler,
		jwtManager:          jwtManager,
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
		companies.DELETE("/:companyId/users/:userId", r.companyHandler.RemoveUserFromCompany)

		// Subscription routes nested under company
		companies.GET("/:companyId/subscription", r.subscriptionHandler.GetSubscription)
		companies.PUT("/:companyId/subscription", r.subscriptionHandler.UpdateSubscription)

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

		// Franchise routes
		companies.POST("/:companyId/franchises", r.franchiseHandler.CreateFranchise)
		companies.GET("/:companyId/franchises", r.franchiseHandler.ListFranchises)

		// Inventory routes
		companies.GET("/:companyId/inventory", r.inventoryHandler.GetCompanyInventory)
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
		franchises.DELETE("/:franchiseId/pricing/:variantId", r.franchiseHandler.DeleteFranchisePricing)
		franchises.POST("/:franchiseId/users", r.franchiseHandler.AddUserToFranchise)
		franchises.DELETE("/:franchiseId/users/:userId", r.franchiseHandler.RemoveUserFromFranchise)
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
