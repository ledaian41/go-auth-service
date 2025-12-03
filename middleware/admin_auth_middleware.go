package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-auth-service/internal/shared/interface"
	"net/http"
)

func AdminAuthMiddleware(authService shared_interface.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}

		mapClaims := claims.(jwt.MapClaims)
		role := mapClaims["role"]
		if role != nil && authService.CheckAdminRole(role.([]interface{})) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "no permission"})
		return
	}
}
