package media

import (
	"github.com/gin-gonic/gin"
	"server/media/route"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/media/list", route.List);
	router.GET("/media/upload", route.Upload);
	router.GET("/media/update", route.Update);
}