package user

import (
	"github.com/gin-gonic/gin"
	"services/server/user/route"
	"services/server/user/repository"
	"services/server/user/service"
	repository2 "services/server/core/repository"
)

type UserModule struct {
	channelRepository *repository.ChannelRepository
	userRepository    *repository.UserRepository
	ChannelService    *service.ChannelService
	UserService       *service.UserService
}

func Init() (*UserModule) {
	channelRepository := repository.NewChannelRepository(repository2.GetRepository())
	userRepository := repository.NewUserRepository(repository2.GetRepository())
	channelService := service.NewChannelService(channelRepository)
	userService := service.NewUserService(channelService, userRepository)

	return &UserModule{
		channelRepository: channelRepository,
		userRepository:    userRepository,
		ChannelService:    channelService,
		UserService:       userService,
	}
}

func (module *UserModule) RegisterRoutes(router *gin.Engine) {
	login := route.NewLoginRoute(module.UserService)
	register := route.NewRegisterRoute(module.UserService)
	logout := route.NewLogoutRoute()

	user := router.Group("/user")

	user.GET("/login", login.Get)
	user.POST("/login", login.Post)
	user.GET("/register", register.Get)
	user.POST("/register", register.Post)
	user.GET("/logout", logout.Get)
}
