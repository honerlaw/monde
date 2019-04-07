package middleware

import (
	"github.com/gin-gonic/gin"
	"server/core/util"
)

func Authorize() (gin.HandlerFunc) {
	return func(c *gin.Context) {
		// if we don't have the payload, just redirect home
		value, exists := c.Get("JWT_PAYLOAD")
		if !exists || value == nil {
			util.Redirect(c, "/")
			return
		}

		c.Next()
	};
}
