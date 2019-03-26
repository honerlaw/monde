package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/media/service"
)

func Update(c *gin.Context) {
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
