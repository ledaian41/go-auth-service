package middleware

import (
	"auth/pkg/site/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SiteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		siteId := c.Param("siteId")
		site := site_service.CheckSiteExists(siteId)
		if site == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Site not found"})
			c.Abort()
			return
		}
		c.Set("site", site)
		c.Next()
	}
}
