package middleware

import (
	"github.com/gin-gonic/gin"
	"go-auth-service/pkg/shared/interface"
	"go-auth-service/pkg/shared/utils"
	"net/http"
)

func AuthMiddleware(authService shared_interface.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "no access token"})
			return
		}

		site, err := shared_utils.ReadSiteContext(c)
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
