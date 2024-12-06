package middleware

import (
	"auth/pkg/auth/service"
	"auth/pkg/auth/utils"
	"auth/pkg/site/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func SiteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		siteId := c.Param("siteId")
		if strings.Trim(siteId, " ") == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "site not found"})
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
		tokenStr, err := utils.GetCookieToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		claims, err := auth_service.ExtractJwtToken(c, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
