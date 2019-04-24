package route

import (
	"github.com/gin-gonic/gin"
	"services/server/core/render"
	"net/http"
	"services/server/core/util"
	"services/server/media/service"
)

type ViewRoute struct {
	mediaService *service.MediaService
	commentService *service.CommentService
}

func NewViewRoute(mediaService *service.MediaService, commentService *service.CommentService) (*ViewRoute) {
	return &ViewRoute{
		mediaService: mediaService,
		commentService: commentService,
	}
}

func (route *ViewRoute) Get(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		util.Redirect(c, "/");
		return;
	}

	props := gin.H{
		"view": nil,
	}

	data, err := route.mediaService.GetByVideoID(id)
	if err != nil {
		props["error"] = err.Error()
	}

	if data != nil {
		if !data.Media.Published {
			util.Redirect(c, "/404");
			return;
		}

		props["view"] = route.mediaService.ConvertSingleMediaInfo(*data, "", "", nil)
	}

	props["comments"], err = route.commentService.GetByMediaID(id)

	render.RenderPage(c, http.StatusOK, props);
}