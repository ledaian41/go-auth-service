package user_handler

import (
	"github.com/gin-gonic/gin"
	"go-auth-service/pkg/shared/interface"
	"net/http"
)

type UserHandler struct {
	userService shared_interface.UserService
}

func NewUserHandler(userService shared_interface.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (handler *UserHandler) GetUserList(c *gin.Context) {
	siteId := c.Param("siteId")
	responses := handler.userService.FindUsersBySite(siteId)
	c.IndentedJSON(http.StatusOK, responses)
}
