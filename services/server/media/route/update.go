package route

import (
	"github.com/gin-gonic/gin"
	"services/server/media/service"
	"services/server/core/util"
)

type UpdateRoute struct {
	mediaService *service.MediaService
}

func NewUpdateRoute(mediaService *service.MediaService) (*UpdateRoute) {
	return &UpdateRoute{
		mediaService: mediaService,
	}
}

func (route *UpdateRoute) Post(c *gin.Context) {
	var req service.UpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Set("error", err.Error())
		util.Redirect(c, "/media/list")
		return
	}

	err := route.mediaService.Update(req)
	if err != nil {
		c.Set("error", err.Error())
		util.Redirect(c, "/media/list")
		return
	}

	util.Redirect(c, "/media/list")
}
