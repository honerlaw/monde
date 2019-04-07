package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"services/server/core/render"
	"services/server/media/service"
)

func Home(c *gin.Context) {
	render.RenderPage(c, http.StatusOK, service.GetHomeMediaResponseProps(c))
}