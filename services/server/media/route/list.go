package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"services/server/media/service"
	"services/server/core/render"
	"services/server/core/util"
	"services/server/user/middleware"
)

func List(c *gin.Context) {
	payload := c.MustGet("JWT_AUTH_PAYLOAD")

	if payload == nil {
		util.Redirect(c, "/")
		return
	}

	// fetch requested media info for given page
	data, err := service.GetMediaDataByUserId(payload.(*middleware.AuthPayload).ID, util.GetSelectPage(c))
	if err != nil {
		render.RenderPage(c, http.StatusInternalServerError, nil)
		return
	}

	props := service.GetListMediaResponseProps(c, data)

	render.RenderPage(c, http.StatusOK, props)
}
