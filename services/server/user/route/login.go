package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/core/render"
	"server/core/util"
)

func Login(c *gin.Context) {
	payload := c.MustGet("JWT_PAYLOAD")

	if payload != nil {
		util.Redirect(c, "/")
		return
	}

	render.RenderPage(c, http.StatusOK, nil)
}
