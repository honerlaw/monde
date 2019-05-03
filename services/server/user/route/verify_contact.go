package route

import (
	"github.com/gin-gonic/gin"
	"encoding/base64"
	"encoding/json"
	"services/server/user/service"
	"services/server/core/util"
	"services/server/core/render"
)

type VerifyContactRoute struct {
	contactService *service.ContactService
}

func NewVerifyContactRoute(contactService *service.ContactService) (*VerifyContactRoute) {
	return &VerifyContactRoute{
		contactService: contactService,
	}
}

func (route *VerifyContactRoute) Get(c *gin.Context) {
	decode, err := base64.StdEncoding.DecodeString(c.Param("data"))
	if err != nil {
		util.RedirectWithError(c, "/500", "failed to verify contact")
		return
	}

	payload := &service.VerifyContactPayload{}
	err = json.Unmarshal(decode, payload)
	if err != nil {
		util.RedirectWithError(c, "/500", "failed to verify contact")
		return
	}

	err = route.contactService.Verify(payload)
	if err != nil {
		util.RedirectWithError(c, "/500", "failed to verify contact")
		return
	}

	render.RenderPage(c, 200, nil)
}
