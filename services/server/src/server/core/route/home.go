package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/core/render"
)

func Home(c *gin.Context) {
	render.RenderPage(c, http.StatusOK, nil)
}