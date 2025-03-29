package auth_utils

import (
	"github.com/gin-gonic/gin"
	"go-auth-service/pkg/shared/dto"
	"net/http"
)

func SetCookieToken(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	site, _ := c.Get("site")
	siteId := site.(*shared_dto.SiteDTO).ID
	c.SetCookie("jwt-"+siteId, token, 3600*24*7, "", "", false, true) // 1 week
}

func GetCookieToken(c *gin.Context) (string, error) {
	site, _ := c.Get("site")
	siteId := site.(*shared_dto.SiteDTO).ID
	return c.Cookie("jwt-" + siteId)
}

func DestroyCookieToken(c *gin.Context) {
	site, _ := c.Get("site")
	siteId := site.(*shared_dto.SiteDTO).ID
	c.SetCookie("jwt-"+siteId, "", 0, "", "", false, true)
}

func ToStringSlice(input []interface{}) []string {
	strSlice := make([]string, len(input))
	for i, v := range input {
		str, ok := v.(string)
		if !ok {
			return []string{}
		}
		strSlice[i] = str
	}
	return strSlice
}
