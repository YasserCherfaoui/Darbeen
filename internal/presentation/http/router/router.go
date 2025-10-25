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
	jwtManager          *security.JWTManager
}

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	companyHandler *handler.CompanyHandler,
	subscriptionHandler *handler.SubscriptionHandler,
	jwtManager *security.JWTManager,
) *Router {
	return &Router{
		authHandler:         authHandler,
		userHandler:         userHandler,
		companyHandler:      companyHandler,
		subscriptionHandler: subscriptionHandler,
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
		companies.GET("/:id", r.companyHandler.GetCompany)
		companies.PUT("/:id", r.companyHandler.UpdateCompany)
		companies.POST("/:id/users", r.companyHandler.AddUserToCompany)
		companies.DELETE("/:id/users/:userId", r.companyHandler.RemoveUserFromCompany)

		// Subscription routes nested under company
		companies.GET("/:id/subscription", r.subscriptionHandler.GetSubscription)
		companies.PUT("/:id/subscription", r.subscriptionHandler.UpdateSubscription)
	}

	// Health check
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
