package route

import (
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	"services/server/core/util"
	"log"
)

const cookieName = "jwt"

type LogoutRoute struct {

}

func NewLogoutRoute() (*LogoutRoute) {
	return &LogoutRoute{}
}

func (route *LogoutRoute) Get(c *gin.Context) {
	cookie, err := c.Request.Cookie(cookieName)

	if err != nil {
		log.Printf("Failed to get cookie information for name: %s, err: %s", cookieName, err)
		return
	}

	cookie.Name = cookieName
	cookie.Path = "/"
	cookie.Value = "invalid"
	cookie.Expires = time.Unix(0, 0)

	http.SetCookie(c.Writer, cookie)

	util.Redirect(c, "/")
}
