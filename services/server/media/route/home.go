package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"services/server/core/render"
	"services/server/media/service"
	"services/server/core/util"
)

type HomeRoute struct {
	mediaService *service.MediaService
}

func NewHomeRoute(mediaService *service.MediaService) (*HomeRoute) {
	return &HomeRoute{
		mediaService: mediaService,
	}
}

func (con *HomeRoute) Get(c *gin.Context) {
	infos, err := con.mediaService.List(util.GetSelectPage(c))

	var props = gin.H{
		"error": err,
		"media": []service.MediaResponse{},
	}

	if err == nil {
		props["media"] = con.mediaService.ConvertMediaData(infos, nil)
	}

	render.RenderPage(c, http.StatusOK, props)
}