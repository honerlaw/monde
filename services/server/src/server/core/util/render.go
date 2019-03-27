package util

import (
	"github.com/gin-gonic/gin"
	"github.com/musawirali/preact-rpc/goclient"
	"net/http"
)

func RenderPage(c *gin.Context, statusCode int, component string, props gin.H) {
	resp, err := goclient.RenderComponent(component, nil, props)
	if err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(statusCode)
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.Write([]byte(resp.Html))
}

func RenderPage500(c *gin.Context) {
	RenderPage(c, http.StatusInternalServerError, "Error500", nil)
}
