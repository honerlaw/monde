package media

import (
	"github.com/gin-gonic/gin"
	"services/server/media/route"
	"services/server/media/repository"
	repository2 "services/server/core/repository"
	"services/server/media/service"
	"services/server/user/middleware"
)

func RegisterRoutes(router *gin.Engine) {
	// repositories
	mediaRepository := repository.NewMediaRepository(repository2.GetRepository())
	hashtagRepository := repository.NewHashtagRepository(repository2.GetRepository())
	commentRepository := repository.NewCommentRepository(repository2.GetRepository())

	// services
	mediaService := service.NewMediaService(mediaRepository, hashtagRepository)
	commentService := service.NewCommentService(commentRepository)

	// routes
	homeRoute := route.NewHomeRoute(mediaService)
	updateRoute := route.NewUpdateRoute(mediaService)
	publishRoute := route.NewPublishRoute(mediaService)
	viewRoute := route.NewViewRoute(mediaService, commentService)
	listRoute := route.NewListRoute(mediaService)
	commentRoute := route.NewCommentRoute(commentService)

	media := router.Group("/media")
	media.Use(middleware.Authorize())

	media.GET("/list", listRoute.Get);
	media.POST("/update", updateRoute.Post);
	media.POST("/publish", publishRoute.Post);
	media.GET("/view/:id", viewRoute.Get);
	media.POST("/comment/:id", commentRoute.Post)

	router.GET("/", homeRoute.Get)
}
