package user

import (
	"go-auth-service/internal/shared/interface"
	"net/http"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, responses)
}
