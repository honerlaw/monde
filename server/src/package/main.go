package main

import (
	"github.com/gin-gonic/gin"
	"github.com/musawirali/preact-rpc/goclient"
	"package/controller"
)

func main() {

	// connect to the react rendering server
	goclient.Connect("tcp", "0.0.0.0:9000")

	router := gin.Default()
	router.Static("/css", "../../../assets/css")
	router.StaticFile("/favicon.ico", "../../../assets/favicon.ico")

	controller.RegisterLoginController(router)

	router.Run("0.0.0.0:8080")
}
