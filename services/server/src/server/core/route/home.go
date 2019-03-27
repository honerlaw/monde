package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/core/util"
)

func Home(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY")

	util.RenderPage(c, http.StatusOK, "HomePage", gin.H{
		"authPayload": payload,
	})
}