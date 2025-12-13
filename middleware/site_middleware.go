package middleware

import (
	"github.com/gin-gonic/gin"
	"go-auth-service/internal/shared"
	"net/http"
	"strings"
)

func SiteMiddleware(siteService shared.SiteService) gin.HandlerFunc {
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
