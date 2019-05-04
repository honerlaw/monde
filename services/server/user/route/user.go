package route

import (
	"services/server/user/service"
	"github.com/gin-gonic/gin"
	"services/server/core/render"
	"net/http"
	"services/server/user/middleware"
)

type UserRoute struct {
	addressService *service.AddressService
}

func NewUserRoute(addressService *service.AddressService) (*UserRoute) {
	return &UserRoute{
		addressService: addressService,
	}
}

func (route *UserRoute) Get(c *gin.Context) {
	payload := c.MustGet("JWT_AUTH_PAYLOAD").(*middleware.AuthPayload)

	addresses := route.addressService.GetAddressesByUserID(payload.ID)

	render.RenderPage(c, http.StatusOK, gin.H{
		"addresses": addresses,
	})
}
