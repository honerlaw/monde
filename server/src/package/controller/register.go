package controller

import (
	"github.com/gin-gonic/gin"
	"package/util"
	"net/http"
	"package/service"
	"package/middleware"
)

func RegisterController(router *gin.Engine) {
	router.GET("/register", renderRegisterPage)
	router.POST("/register", handleRegister)
}

func renderRegisterPage(c *gin.Context) {
	util.RenderPage(c, http.StatusOK, "RegisterPage", nil)
}

func handleRegister(c *gin.Context) {
	var req service.CreateUserRequest

	if err := c.ShouldBind(&req); err != nil {
		util.RenderPage(c, http.StatusBadRequest, "RegisterPage", gin.H{
			"error": "all fields are required",
		})
		return
	}

	_, err := service.CreateUser(req)

	if err != nil {
		util.RenderPage(c, http.StatusBadRequest, "RegisterPage", gin.H{
			"error": err.Error(),
		})
		return
	}

	middleware.AuthMiddleware.LoginHandler(c)
}
