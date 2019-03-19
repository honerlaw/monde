package controller

import (
	"github.com/gin-gonic/gin"
	"package/util"
	"net/http"
)

func RegisterController(router *gin.Engine) {
	router.GET("/register", renderRegisterPage)
}

func renderRegisterPage(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY")

	if payload != nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}

	util.RenderPage(c, http.StatusOK, "RegisterPage", nil)
}
