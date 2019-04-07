package media

import (
	"github.com/gin-gonic/gin"
	"services/server/media/route"
	"services/server/user/middleware"
)

func RegisterRoutes(router *gin.Engine) {
	media := router.Group("/media")
	media.Use(middleware.Authorize())

	media.GET("/list", route.List);
	media.POST("/update", route.Update);
	media.POST("/publish", route.Publish);
	media.GET("/view/:id", route.View);
	router.GET("/", route.Home)
}
