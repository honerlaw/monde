package media

import (
	"github.com/gin-gonic/gin"
	"package/util"
	"net/http"
)

func List(c *gin.Context) {
	payload := c.MustGet("JWT_IDENTITY")
	props := gin.H{
		"authPayload": payload,
	}

	if payload != nil {
		// 1. fetch all mediainfo for the user
		// 2. fetch the job id and see if it finished

	}

	util.RenderPage(c, http.StatusOK, "UploadListPage", props)
}
