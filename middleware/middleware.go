package middleware

import (
	auth_service "auth/pkg/auth/service"
	site_service "auth/pkg/site/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func SiteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		siteId := c.Param("siteId")
		if strings.Trim(siteId, " ") == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "siteId not found"})
			return
		}

		site := site_service.CheckSiteExists(siteId)
		if site == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "site not found"})
			return
		}

		c.Set("site", site)
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "no access token"})
			return
		}

		claims, err := auth_service.ValidateAccessToken(c, accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
