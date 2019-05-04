package route

import (
	"github.com/gin-gonic/gin"
	"encoding/base64"
	"encoding/json"
	"services/server/user/service"
	"services/server/core/util"
	"services/server/core/render"
	"net/http"
)

type ContactVerifyRoute struct {
	contactService *service.ContactService
}

func NewContactVerifyRoute(contactService *service.ContactService) (*ContactVerifyRoute) {
	return &ContactVerifyRoute{
		contactService: contactService,
	}
}

func (route *ContactVerifyRoute) Get(c *gin.Context) {
	decode, err := base64.StdEncoding.DecodeString(c.Param("data"))
	if err != nil {
		util.RedirectWithError(c, "/500", http.StatusInternalServerError, "failed to verify contact")
		return
	}

	payload := &service.ContactVerifyPayload{}
	err = json.Unmarshal(decode, payload)
	if err != nil {
		util.RedirectWithError(c, "/500", http.StatusInternalServerError, "failed to verify contact")
		return
	}

	err = route.contactService.Verify(payload)
	if err != nil {
		util.RedirectWithError(c, "/500", http.StatusInternalServerError, "failed to verify contact")
		return
	}

	render.RenderPage(c, 200, nil)
}
