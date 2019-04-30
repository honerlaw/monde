package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CoreModule struct {
}

func Init() (*CoreModule) {
	return &CoreModule{}
}

func (module *CoreModule) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})
}
