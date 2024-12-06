package utils

import (
	"auth/pkg/site/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetCookieToken(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	site, _ := c.Get("site")
	siteId := site.(*site_model.Site).ID
	c.SetCookie("jwt-"+siteId, token, 3600*24*7, "", "", false, true) // 1 week
}

func GetCookieToken(c *gin.Context) (string, error) {
	site, _ := c.Get("site")
	siteId := site.(*site_model.Site).ID
	return c.Cookie("jwt-" + siteId)
}

func DestroyCookieToken(c *gin.Context) {
	site, _ := c.Get("site")
	siteId := site.(*site_model.Site).ID
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
