package route

import (
	"services/server/user/service"
	"github.com/gin-gonic/gin"
	"services/server/core/render"
	"net/http"
	"services/server/user/middleware"
)

type UserRoute struct {
	contactService *service.ContactService
	addressService *service.AddressService
}

func NewUserRoute(contactService *service.ContactService, addressService *service.AddressService) (*UserRoute) {
	return &UserRoute{
		contactService: contactService,
		addressService: addressService,
	}
}

func (route *UserRoute) Get(c *gin.Context) {
	payload := c.MustGet("JWT_AUTH_PAYLOAD").(*middleware.AuthPayload)

	contacts := route.contactService.GetContactDataByUserID(payload.ID)
	addresses := route.addressService.GetAddressesByUserID(payload.ID)

	render.RenderPage(c, http.StatusOK, gin.H{
		"contacts": contacts,
		"addresses": addresses,
	})
}
