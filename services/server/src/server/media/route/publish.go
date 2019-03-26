package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/media/service"
)

func Publish(c *gin.Context) {
	var req service.PublishRequest
	if err := c.ShouldBind(&req); err != nil {
		// @todo show error
		c.Redirect(http.StatusFound, "/media/list")
		c.Abort()
		return
	}

	err := service.TogglePublish(req)
	if err != nil {
		// @todo show error
		c.Redirect(http.StatusFound, "/media/list")
		c.Abort()
		return
	}

	c.Redirect(http.StatusFound, "/media/list")
	c.Abort()
}