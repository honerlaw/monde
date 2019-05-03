package user

import (
	"github.com/gin-gonic/gin"
	"services/server/user/route"
	"services/server/user/repository"
	"services/server/user/service"
	repository2 "services/server/core/repository"
)

type UserModule struct {
	ContactRepository *repository.ContactRepository
	channelRepository *repository.ChannelRepository
	userRepository    *repository.UserRepository
	ContactService    *service.ContactService
	ChannelService    *service.ChannelService
	UserService       *service.UserService
}

func Init() (*UserModule) {
	contactRepository := repository.NewContactRepository(repository2.GetRepository())
	channelRepository := repository.NewChannelRepository(repository2.GetRepository())
	userRepository := repository.NewUserRepository(repository2.GetRepository())
	contactService := service.NewContactService(contactRepository)
	channelService := service.NewChannelService(channelRepository)
	userService := service.NewUserService(contactService, channelService, userRepository)

	return &UserModule{
		ContactRepository: contactRepository,
		channelRepository: channelRepository,
		userRepository:    userRepository,
		ContactService:    contactService,
		ChannelService:    channelService,
		UserService:       userService,
	}
}

func (module *UserModule) RegisterRoutes(router *gin.Engine) {
	login := route.NewLoginRoute(module.UserService)
	register := route.NewRegisterRoute(module.UserService)
	logout := route.NewLogoutRoute()
	verifyContact := route.NewVerifyContactRoute(module.ContactService)

	user := router.Group("/user")
	user.GET("/login", login.Get)
	user.POST("/login", login.Post)
	user.GET("/register", register.Get)
	user.POST("/register", register.Post)
	user.GET("/logout", logout.Get)

	contact := router.Group("/contact")
	contact.GET("/verify/:data", verifyContact.Get)
}
