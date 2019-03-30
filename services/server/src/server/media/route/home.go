package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/core/render"
	"server/media/service"
)

func Home(c *gin.Context) {
	render.RenderPage(c, http.StatusOK, service.GetHomeMediaResponseProps(c))
}