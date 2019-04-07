package route

import (
	"github.com/gin-gonic/gin"
	"services/server/media/service"
	"services/server/core/util"
)

func Publish(c *gin.Context) {
	var req service.PublishRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Set("error", err)
		util.Redirect(c, "/media/list")
		return
	}

	err := service.TogglePublish(req)
	if err != nil {
		c.Set("error", err)
		util.Redirect(c, "/media/list")
		return
	}

	util.Redirect(c, "/media/list")
}