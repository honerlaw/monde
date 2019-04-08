package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"services/server/core/render"
	"services/server/core/util"
)

func Login(c *gin.Context) {
	payload := c.MustGet("JWT_AUTH_PAYLOAD")

	if payload != nil {
		util.Redirect(c, "/")
		return
	}

	render.RenderPage(c, http.StatusOK, nil)
}
