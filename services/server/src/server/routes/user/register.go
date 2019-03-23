package user

import (
	"github.com/gin-gonic/gin"
	"server/util"
	"net/http"
)

func Register(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY")

	if payload != nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}

	util.RenderPage(c, http.StatusOK, "RegisterPage", nil)
}
