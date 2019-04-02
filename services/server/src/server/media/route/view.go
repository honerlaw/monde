package route

import (
	"github.com/gin-gonic/gin"
	"server/core/render"
	"net/http"
	"server/core/util"
	"server/media/service"
)

func View(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		util.Redirect(c, "/");
		return;
	}

	props := gin.H{
		"view": nil,
	}

	info, err := service.GetMediaInfoByVideoID(id)
	if err != nil {
		props["error"] = err.Error()
	}

	if info != nil {
		props["view"] = service.ConvertSingleMediaInfo(*info, "", "", nil)
	}

	render.RenderPage(c, http.StatusOK, props);
}