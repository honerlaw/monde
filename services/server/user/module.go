package user

import (
	"github.com/gin-gonic/gin"
	"services/server/user/route"
	"services/server/user/repository"
	"services/server/user/service"
	repository2 "services/server/core/repository"
	"services/server/user/middleware"
)

type UserModule struct {
	addressRepository *repository.AddressRepository
	contactRepository *repository.ContactRepository
	channelRepository *repository.ChannelRepository
	userRepository    *repository.UserRepository
	addressService    *service.AddressService
	ContactService    *service.ContactService
	ChannelService    *service.ChannelService
	UserService       *service.UserService
}

func Init() (*UserModule) {
	addressRepository := repository.NewAddressRepository(repository2.GetRepository())
	contactRepository := repository.NewContactRepository(repository2.GetRepository())
	channelRepository := repository.NewChannelRepository(repository2.GetRepository())
	userRepository := repository.NewUserRepository(repository2.GetRepository())
	addressService := service.NewAddressService(addressRepository)
	contactService := service.NewContactService(contactRepository)
	channelService := service.NewChannelService(channelRepository)
	userService := service.NewUserService(contactService, channelService, userRepository)

	return &UserModule{
		contactRepository: contactRepository,
		channelRepository: channelRepository,
		userRepository:    userRepository,
		addressService:    addressService,
		ContactService:    contactService,
		ChannelService:    channelService,
		UserService:       userService,
	}
}

func (module *UserModule) RegisterRoutes(router *gin.Engine) {
	login := route.NewLoginRoute(module.UserService)
	register := route.NewRegisterRoute(module.UserService)
	logout := route.NewLogoutRoute()
	contactVerify := route.NewContactVerifyRoute(module.ContactService)
	address := route.NewAddressCreateRoute(module.addressService)
	user := route.NewUserRoute(module.ContactService, module.addressService)

	userGroup := router.Group("/user")
	userGroup.GET("/login", login.Get)
	userGroup.POST("/login", login.Post)
	userGroup.GET("/register", register.Get)
	userGroup.POST("/register", register.Post)
	userGroup.GET("/logout", logout.Get)
	userGroup.GET("", middleware.Authorize(), user.Get)

	contact := router.Group("/contact")
	contact.GET("/verify/:data", middleware.Authorize(), contactVerify.Get)

	addressGroup := router.Group("/address")
	addressGroup.POST("", middleware.Authorize(), address.Post)
}
