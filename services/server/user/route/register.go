package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"services/server/user/service"
	"services/server/user/middleware"
	"services/server/core/render"
	"services/server/core/util"
)

func Register(c *gin.Context) {
	payload := c.MustGet("JWT_AUTH_PAYLOAD")

	if payload != nil {
		util.Redirect(c, "/")
		return
	}

	render.RenderPage(c, http.StatusOK, nil)
}

func RegisterPost(c *gin.Context) {
	var req service.CreateRequest

	if err := c.ShouldBind(&req); err != nil {
		render.RenderPage(c, http.StatusBadRequest, gin.H{
			"usernname": req.Username,
			"error": "all fields are required",
		})
		return
	}

	_, err := service.Create(req)

	if err != nil {
		render.RenderPage(c, http.StatusBadRequest, gin.H{
			"usernname": req.Username,
			"error": err.Error(),
		})
		return
	}

	middleware.GetJWTAuth().LoginHandler(c)
}