package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authorize() (gin.HandlerFunc) {
	return func(c *gin.Context) {
		// if we don't have the payload, just redirect home
		value, exists := c.Get("JWT_IDENTITY")
		if !exists || value == nil {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		c.Next()
	};
}
