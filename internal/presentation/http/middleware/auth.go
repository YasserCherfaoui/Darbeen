package middleware

import (
	"strings"

	"github.com/YasserCherfaoui/darween/internal/infrastructure/security"
	"github.com/YasserCherfaoui/darween/internal/presentation/response"
	"github.com/YasserCherfaoui/darween/pkg/errors"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "userID"
	UserEmailKey        = "userEmail"
)

func AuthMiddleware(jwtManager *security.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Error(c, errors.NewUnauthorizedError("authorization header is required"))
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			response.Error(c, errors.NewUnauthorizedError("invalid authorization format"))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, BearerPrefix)
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			response.Error(c, errors.NewUnauthorizedError("invalid or expired token"))
			c.Abort()
			return
		}

		// Set user information in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, errors.NewUnauthorizedError("user not authenticated")
	}
	return userID.(uint), nil
}
