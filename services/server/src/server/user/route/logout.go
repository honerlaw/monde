package route

import (
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	"server/user/middleware"
	"server/core/util"
	"log"
)

func Logout(c *gin.Context) {
	cookieName := middleware.GetJWTAuth().CookieName
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
