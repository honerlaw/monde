package route

import (
	"github.com/gin-gonic/gin"
	"services/server/core/render"
	"net/http"
	"services/server/core/util"
	"services/server/media/service"
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

	data, err := service.GetMediaDataByVideoID(id)
	if err != nil {
		props["error"] = err.Error()
	}

	if data != nil {
		props["view"] = service.ConvertSingleMediaInfo(*data, "", "", nil)
	}

	render.RenderPage(c, http.StatusOK, props);
}