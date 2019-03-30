package route

import (
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	"server/user/middleware"
)

func Logout(c *gin.Context) {
	cookieName := middleware.GetJWTAuth().CookieName
	cookie, err := c.Request.Cookie(cookieName)

	if err != nil {
		panic(err)
	}

	cookie.Name = cookieName
	cookie.Path = "/"
	cookie.Value = "invalid"
	cookie.Expires = time.Unix(0, 0)

	http.SetCookie(c.Writer, cookie)

	c.Redirect(http.StatusFound, "/")
	c.Abort()
}
