package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"services/server/core/render"
	"services/server/core/util"
	"services/server/user/service"
	"services/server/user/middleware"
)

type LoginRoute struct {
	userService *service.UserService
}

func NewLoginRoute(userService *service.UserService) (*LoginRoute) {
	return &LoginRoute{
		userService: userService,
	}
}

func (route *LoginRoute) Get(c *gin.Context) {
	payload := c.MustGet("JWT_AUTH_PAYLOAD")

	if payload != nil {
		util.Redirect(c, "/")
		return
	}

	render.RenderPage(c, http.StatusOK, nil)
}

func (route *LoginRoute) Post(c *gin.Context) {
	middleware.GetJWTAuth().LoginHandler(c)
}