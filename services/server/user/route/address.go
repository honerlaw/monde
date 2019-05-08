package route

import (
	"services/server/user/service"
	"github.com/gin-gonic/gin"
	"services/server/core/util"
	"net/http"
	"services/server/user/middleware"
)

type AddressRoute struct {
	addressService *service.AddressService
}

func NewAddressCreateRoute(addressService *service.AddressService) (*AddressRoute) {
	return &AddressRoute{
		addressService: addressService,
	}
}

func (route *AddressRoute) Post(c *gin.Context) {
	payload := c.MustGet("JWT_AUTH_PAYLOAD").(*middleware.AuthPayload)

	var req service.AddressCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		util.RedirectWithError(c, "/user", http.StatusBadRequest, "all fields are required")
		return
	}

	_, err := route.addressService.CreateOrUpdate(payload.ID, &req)
	if err != nil {
		util.RedirectWithError(c, "/user", http.StatusInternalServerError, err.Error())
		return
	}

	util.Redirect(c, "/user")
}

func (route *AddressRoute) Delete(c *gin.Context) {
	payload := c.MustGet("JWT_AUTH_PAYLOAD").(*middleware.AuthPayload)
	id := c.Param("id")

	err := route.addressService.Delete(payload.ID, id)
	if err != nil {
		util.RedirectWithError(c, "/user", http.StatusInternalServerError, err.Error())
		return
	}

	util.Redirect(c, "/user")
}