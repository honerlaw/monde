package controller

import (
	"github.com/gin-gonic/gin"
	"package/util"
	"net/http"
)

func LoginController(router *gin.Engine) {
	router.GET("/login", renderLoginPage)
}

func renderLoginPage(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY")

	if payload != nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}

	util.RenderPage(c, http.StatusOK, "LoginPage", nil)
}
