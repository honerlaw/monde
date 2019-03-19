package routes

import (
	"github.com/gin-gonic/gin"
	"package/util"
	"net/http"
)

func Home(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY")

	util.RenderPage(c, http.StatusOK, "HomePage", gin.H{
		"authPayload": payload,
	})
}