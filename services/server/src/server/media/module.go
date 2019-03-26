package media

import (
	"github.com/gin-gonic/gin"
	"server/media/route"
	"server/middleware/auth"
)

func RegisterRoutes(router *gin.Engine) {
	media := router.Group("/media")
	media.Use(auth.Authorize())

	media.GET("/list", route.List);
	media.GET("/upload", route.Upload);
	media.POST("/update", route.Update);
	media.POST("/publish", route.Publish);
}
