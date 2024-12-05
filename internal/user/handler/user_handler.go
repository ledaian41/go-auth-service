package user_handler

import (
	"auth/internal/user/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserList(c *gin.Context) {
	responses := make([]user_model.UserResponse, len(user_model.UserList))
	for i, user := range user_model.UserList {
		responses[i] = user.Response()
	}
	c.IndentedJSON(http.StatusOK, responses)
}
