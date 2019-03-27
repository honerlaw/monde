package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/core/util"
	"server/user/service"
	"server/user/middleware"
)

func Register(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY")

	if payload != nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}

	util.RenderPage(c, http.StatusOK, "RegisterPage", nil)
}

func RegisterPost(c *gin.Context) {
	var req service.CreateRequest

	if err := c.ShouldBind(&req); err != nil {
		util.RenderPage(c, http.StatusBadRequest, "RegisterPage", gin.H{
			"usernname": req.Username,
			"error": "all fields are required",
		})
		return
	}

	_, err := service.Create(req)

	if err != nil {
		util.RenderPage(c, http.StatusBadRequest, "RegisterPage", gin.H{
			"usernname": req.Username,
			"error": err.Error(),
		})
		return
	}

	middleware.GetJWTAuth().LoginHandler(c)
}