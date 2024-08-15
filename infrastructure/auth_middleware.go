package infrastructure

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddlewareService interface {
	AuthMiddleware() gin.HandlerFunc
	AdminMiddleware() gin.HandlerFunc
}

type AuthService struct {
	jwtService JWTService
}

func NewAuthService(jwtService JWTService) AuthMiddlewareService {
	return &AuthService{jwtService: jwtService}
}




func (am *AuthService) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization format must be Bearer {token}"})
			return
		}

		tokenString := parts[1]
		claims, err := am.jwtService.ValidateToken(tokenString)
		if err.ErrCode != 0  {
			c.AbortWithStatusJSON(err.ErrCode, gin.H{"message": err.ErrMessage})
			return
		}

		c.Set("userId", claims["userId"])
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])
		c.Next()
	}
}


func (am *AuthService) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Admins only"})
			c.Abort()
			return
		}
		c.Next()
	}
}