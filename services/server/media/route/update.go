package route

import (
	"github.com/gin-gonic/gin"
	"server/media/service"
	"server/core/util"
)

func Update(c *gin.Context) {
	var req service.UpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Set("error", err.Error())
		util.Redirect(c, "/media/list")
		return
	}

	err := service.Update(req)
	if err != nil {
		c.Set("error", err.Error())
		util.Redirect(c, "/media/list")
		return
	}

	util.Redirect(c, "/media/list")
}
