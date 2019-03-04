package controller

import (
	"github.com/gin-gonic/gin"
	"package/util"
	"net/http"
	"package/middleware"
	"time"
)

func LoginController(router *gin.Engine) {
	router.GET("/login", renderLoginPage)
	router.POST("/login", middleware.AuthMiddleware.LoginHandler)

	router.GET("/logout", handleLogout)
}

func renderLoginPage(c *gin.Context) {
	payload := middleware.GetAuthPayload(c)

	if payload != nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}

	util.RenderPage(c, http.StatusOK, "LoginPage", nil)
}

func handleLogout(c *gin.Context) {
	cookieName := middleware.AuthMiddleware.CookieName
	cookie, err := c.Request.Cookie(cookieName)

	if err != nil {
		panic(err)
	}

	cookie.Value = "invalid"
	cookie.Expires = time.Unix(0, 0)

	http.SetCookie(c.Writer, cookie)

	c.Redirect(http.StatusFound, "/")
	c.Abort()
}
