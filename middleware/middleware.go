package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-auth-service/pkg/shared/interface"
	"net/http"
	"strings"
)

func SiteMiddleware(siteService shared_interface.SiteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		siteId := c.Param("siteId")
		if strings.Trim(siteId, " ") == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "siteId not found"})
			return
		}

		site := siteService.CheckSiteExists(siteId)
		if site == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "site not found"})
			return
		}

		c.Set("site", site)
		c.Next()
	}
}

func AuthMiddleware(authService shared_interface.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "no access token"})
			return
		}

		claims, err := authService.ValidateAccessToken(c, accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

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
