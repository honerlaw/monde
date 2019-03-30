package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})
}