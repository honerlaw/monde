package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"services/server/user/service"
	"services/server/user/middleware"
	"services/server/core/render"
	"services/server/core/util"
)

type RegisterRoute struct {
	userService *service.UserService
}

func NewRegisterRoute(userService *service.UserService) (*RegisterRoute) {
	return &RegisterRoute{
		userService: userService,
	}
}

func (route *RegisterRoute) Get(c *gin.Context) {
	payload := c.MustGet("JWT_AUTH_PAYLOAD")

	if payload != nil {
		util.Redirect(c, "/")
		return
	}

	render.RenderPage(c, http.StatusOK, nil)
}

func (route *RegisterRoute) Post(c *gin.Context) {
	var req service.CreateRequest

	if err := c.ShouldBind(&req); err != nil {
		render.RenderPage(c, http.StatusBadRequest, gin.H{
			"email": req.Email,
			"error": "all fields are required",
		})
		return
	}

	_, err := route.userService.Create(req)

	if err != nil {
		render.RenderPage(c, http.StatusBadRequest, gin.H{
			"email": req.Email,
			"error": err.Error(),
		})
		return
	}

	middleware.GetJWTAuth().LoginHandler(c)
}
