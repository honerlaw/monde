package controller

import (
	"github.com/gin-gonic/gin"
	"package/middleware"
	"package/util"
	"net/http"
)

func HomeController(router *gin.Engine) {
	router.GET("/", renderHomePage)
}

func renderHomePage(c *gin.Context) {
	payload := middleware.GetAuthPayload(c);

	util.RenderPage(c, http.StatusOK, "HomePage", gin.H{
		"authPayload": payload,
	})
}