package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/core/render"
)

func Login(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY")

	if payload != nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}

	render.RenderPage(c, http.StatusOK, nil)
}
