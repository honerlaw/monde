package route

import (
	"github.com/gin-gonic/gin"
	"services/server/media/service"
	"services/server/core/util"
)

type PublishRoute struct {
	mediaService *service.MediaService
}

func NewPublishRoute(mediaService *service.MediaService) (*PublishRoute) {
	return &PublishRoute{
		mediaService: mediaService,
	}
}

func (route *PublishRoute) Post(c *gin.Context) {
	var req service.PublishRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Set("error", err)
		util.Redirect(c, "/media/list")
		return
	}

	err := route.mediaService.TogglePublish(req)
	if err != nil {
		c.Set("error", err)
		util.Redirect(c, "/media/list")
		return
	}

	util.Redirect(c, "/media/list")
}