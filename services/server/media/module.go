package media

import (
	"github.com/gin-gonic/gin"
	"services/server/media/route"
	"services/server/user/middleware"
	"services/server/media/repository"
	"services/server/media/service"
	repository2 "services/server/core/repository"
	service2 "services/server/user/service"
)

type MediaModule struct {
	channelService *service2.ChannelService

	mediaRepository   *repository.MediaRepository
	hashtagRepository *repository.HashtagRepository
	commentRepository *repository.CommentRepository
	MediaService      *service.MediaService
	CommentService    *service.CommentService
}

func Init(channelService *service2.ChannelService) (*MediaModule) {
	// repositories
	mediaRepository := repository.NewMediaRepository(repository2.GetRepository())
	hashtagRepository := repository.NewHashtagRepository(repository2.GetRepository())
	commentRepository := repository.NewCommentRepository(repository2.GetRepository())

	// services
	mediaService := service.NewMediaService(mediaRepository, hashtagRepository)
	commentService := service.NewCommentService(commentRepository)

	return &MediaModule{
		// external services
		channelService: channelService,

		// repositories
		mediaRepository:   mediaRepository,
		hashtagRepository: hashtagRepository,
		commentRepository: commentRepository,

		// services
		MediaService:   mediaService,
		CommentService: commentService,
	}
}

func (module *MediaModule) RegisterRoutes(router *gin.Engine) {
	homeRoute := route.NewHomeRoute(module.MediaService)
	updateRoute := route.NewUpdateRoute(module.MediaService)
	publishRoute := route.NewPublishRoute(module.MediaService)
	viewRoute := route.NewViewRoute(module.MediaService, module.CommentService)
	listRoute := route.NewListRoute(module.MediaService, module.channelService)
	commentRoute := route.NewCommentRoute(module.CommentService)

	media := router.Group("/media")
	media.Use(middleware.Authorize())

	media.GET("/list", listRoute.Get);
	media.POST("/update", updateRoute.Post);
	media.POST("/publish", publishRoute.Post);
	media.GET("/view/:id", viewRoute.Get);
	media.POST("/comment/:id", commentRoute.Post)

	router.GET("/", homeRoute.Get)
}
