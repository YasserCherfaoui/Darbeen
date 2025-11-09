package main

import (
	"log"

	"github.com/YasserCherfaoui/darween/internal/application/auth"
	"github.com/YasserCherfaoui/darween/internal/application/company"
	"github.com/YasserCherfaoui/darween/internal/application/franchise"
	"github.com/YasserCherfaoui/darween/internal/application/inventory"
	"github.com/YasserCherfaoui/darween/internal/application/pos"
	"github.com/YasserCherfaoui/darween/internal/application/product"
	"github.com/YasserCherfaoui/darween/internal/application/subscription"
	"github.com/YasserCherfaoui/darween/internal/application/supplier"
	"github.com/YasserCherfaoui/darween/internal/application/user"
	"github.com/YasserCherfaoui/darween/internal/infrastructure/persistence/migrations"
	"github.com/YasserCherfaoui/darween/internal/infrastructure/persistence/postgres"
	"github.com/YasserCherfaoui/darween/internal/infrastructure/security"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/handler"
	"github.com/YasserCherfaoui/darween/internal/presentation/http/router"
	"github.com/YasserCherfaoui/darween/pkg/config"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database
	db, err := postgres.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := migrations.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	companyRepo := postgres.NewCompanyRepository(db)
	subscriptionRepo := postgres.NewSubscriptionRepository(db)
	productRepo := postgres.NewProductRepository(db)
	supplierRepo := postgres.NewSupplierRepository(db)
	franchiseRepo := postgres.NewFranchiseRepository(db)
	inventoryRepo := postgres.NewInventoryRepository(db)
	
	// Initialize POS repositories
	customerRepo := postgres.NewCustomerRepository(db)
	saleRepo := postgres.NewSaleRepository(db)
	saleItemRepo := postgres.NewSaleItemRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	cashDrawerRepo := postgres.NewCashDrawerRepository(db)
	cashDrawerTransactionRepo := postgres.NewCashDrawerTransactionRepository(db)
	refundRepo := postgres.NewRefundRepository(db)

	// Initialize JWT manager
	jwtManager := security.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)

	// Initialize services
	authService := auth.NewService(userRepo, jwtManager)
	userService := user.NewService(userRepo)
	companyService := company.NewService(companyRepo, userRepo, subscriptionRepo)
	subscriptionService := subscription.NewService(subscriptionRepo, userRepo)
	productService := product.NewService(productRepo, userRepo, supplierRepo)
	supplierService := supplier.NewService(supplierRepo, userRepo)
	inventoryService := inventory.NewService(inventoryRepo, companyRepo, franchiseRepo, userRepo, productRepo)
	franchiseService := franchise.NewService(franchiseRepo, inventoryRepo, companyRepo, userRepo, productRepo)
	posService := pos.NewService(customerRepo, saleRepo, saleItemRepo, paymentRepo, cashDrawerRepo, cashDrawerTransactionRepo, refundRepo, userRepo, inventoryRepo, inventoryRepo, productRepo, db)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	companyHandler := handler.NewCompanyHandler(companyService)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	productHandler := handler.NewProductHandler(productService)
	supplierHandler := handler.NewSupplierHandler(supplierService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	franchiseHandler := handler.NewFranchiseHandler(franchiseService)
	posHandler := handler.NewPOSHandler(posService)

	// Initialize router
	r := router.NewRouter(authHandler, userHandler, companyHandler, subscriptionHandler, productHandler, supplierHandler, inventoryHandler, franchiseHandler, posHandler, jwtManager)

	// Create Gin engine
	engine := gin.Default()

	// Setup routes
	r.SetupRoutes(engine)

	// Start server
	serverAddr := ":" + cfg.Server.Port
	log.Printf("Starting server on %s", serverAddr)
	if err := engine.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
