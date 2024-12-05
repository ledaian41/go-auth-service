package auth_utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetCookieToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	token := "test-jwt-token"
	SetCookieToken(c, token)

	resp := w.Result()
	cookies := resp.Cookies()

	assert.Len(t, cookies, 1, "Should have exactly one cookie")
	assert.Equal(t, "jwt", cookies[0].Name, "Cookie name should be 'jwt'")
	assert.Equal(t, token, cookies[0].Value, "Cookie value should match the token")
	assert.Equal(t, http.SameSiteLaxMode, cookies[0].SameSite, "Cookie SameSite attribute should be Lax")
	assert.Equal(t, 3600*24*7, cookies[0].MaxAge, "Cookie MaxAge should be 7 days")
	assert.True(t, cookies[0].HttpOnly, "Cookie should be HttpOnly")
	assert.False(t, cookies[0].Secure, "Cookie should not be Secure")
}
