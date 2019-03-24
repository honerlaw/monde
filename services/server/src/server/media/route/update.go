package route

import (
	"github.com/gin-gonic/gin"
	"server/middleware/auth"
	"net/http"
	"server/media/service"
)

func Update(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY").(*auth.AuthPayload)
	if payload == nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}

	var req service.UpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		// @todo show error
		c.Redirect(http.StatusFound, "/media/list")
		c.Abort()
		return
	}

	err := service.Update(req)
	if err != nil {
		// @todo show error
		c.Redirect(http.StatusFound, "/media/list")
		c.Abort()
		return
	}

	c.Redirect(http.StatusFound, "/media/list")
	c.Abort()
}
