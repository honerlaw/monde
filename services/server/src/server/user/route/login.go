package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/core/util"
)

func Login(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY")

	if payload != nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}

	util.RenderPage(c, http.StatusOK, "LoginPage", nil)
}
