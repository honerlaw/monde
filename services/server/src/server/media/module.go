package media

import (
	"github.com/gin-gonic/gin"
	"server/media/route"
	"server/user/middleware"
	"server/media/model"
	"server/core/repository"
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

func Migrate() {
	(&model.MediaInfo{}).Migrate(repository.DB, repository.MigrateModel)
	(&model.Media{}).Migrate(repository.DB, repository.MigrateModel)
	(&model.Track{}).Migrate(repository.DB, repository.MigrateModel)
	(&model.Hashtag{}).Migrate(repository.DB, repository.MigrateModel)
}
