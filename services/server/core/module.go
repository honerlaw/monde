package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"services/server/core/service"
)

type CoreModule struct {
	CountryService *service.CountryService
}

func Init(countryDataPath string) (*CoreModule) {
	return &CoreModule{
		CountryService: service.NewCountryService(countryDataPath),
	}
}

func (module *CoreModule) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})
}
