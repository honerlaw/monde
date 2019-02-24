package controller

import (
	"github.com/gin-gonic/gin"
	"package/util"
	"net/http"
)

type LoginRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func RegisterLoginController(router *gin.Engine) {
	router.GET("/login", renderLoginPage)
	router.POST("/login", handleLogin)
}

func renderLoginPage(c *gin.Context) {
	util.RenderPage(c, http.StatusOK, "LoginPage", nil)
}

func handleLogin(c *gin.Context) {
	var req LoginRequest
	c.Bind(&req)

	util.RenderPage(c, http.StatusOK, "LoginPage", nil)
}
