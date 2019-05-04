package route

import (
	"services/server/user/service"
	"github.com/gin-gonic/gin"
	"services/server/user/middleware"
	"services/server/core/util"
	"net/http"
)

type AddressRoute struct {
	addressService *service.AddressService
}

func NewAddressCreateRoute(addressService *service.AddressService) (*AddressRoute) {
	return &AddressRoute{
		addressService: addressService,
	}
}

func (route *AddressRoute) Get(c *gin.Context) {
	payload := c.MustGet("JWT_AUTH_PAYLOAD")

	route.addressService.FindByUserID(payload.(*middleware.AuthPayload).ID)
	// @todo should render list of all addresses with create button / modify button
}

func (route *AddressRoute) Put(c *gin.Context) {
	// @todo should simply update everything
}

func (route *AddressRoute) Post(c *gin.Context) {
	payload := c.MustGet("JWT_AUTH_PAYLOAD")

	var req service.AddressCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		util.RedirectWithError(c, "/address", http.StatusBadRequest, "all fields are required")
		return
	}

	_, err := route.addressService.Create(payload.(*middleware.AuthPayload).ID, &req)
	if err != nil {
		util.RedirectWithError(c, "/address", http.StatusInternalServerError, err.Error())
		return
	}

	util.Redirect(c, "/address")
}
