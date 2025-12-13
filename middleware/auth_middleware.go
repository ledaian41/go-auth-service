package middleware

import (
	"go-auth-service/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService shared.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "no access token"})
			return
		}

		site, err := shared.ReadSiteContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		claims, err := authService.ValidateAccessToken(site, accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
