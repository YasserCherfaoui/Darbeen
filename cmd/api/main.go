package main

import (
	"log"

	"github.com/YasserCherfaoui/darween/internal/application/auth"
	"github.com/YasserCherfaoui/darween/internal/application/company"
	"github.com/YasserCherfaoui/darween/internal/application/product"
	"github.com/YasserCherfaoui/darween/internal/application/subscription"
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
	if cfg.Server.GinMode != "debug" {
		if err := migrations.AutoMigrate(db); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	companyRepo := postgres.NewCompanyRepository(db)
	subscriptionRepo := postgres.NewSubscriptionRepository(db)
	productRepo := postgres.NewProductRepository(db)

	// Initialize JWT manager
	jwtManager := security.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)

	// Initialize services
	authService := auth.NewService(userRepo, jwtManager)
	userService := user.NewService(userRepo)
	companyService := company.NewService(companyRepo, userRepo, subscriptionRepo)
	subscriptionService := subscription.NewService(subscriptionRepo, userRepo)
	productService := product.NewService(productRepo, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	companyHandler := handler.NewCompanyHandler(companyService)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	productHandler := handler.NewProductHandler(productService)

	// Initialize router
	r := router.NewRouter(authHandler, userHandler, companyHandler, subscriptionHandler, productHandler, jwtManager)

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
