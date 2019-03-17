package main

import (
	"github.com/gin-gonic/gin"
	"github.com/musawirali/preact-rpc/goclient"
	"package/controller"
	"github.com/joho/godotenv"
	"os"
	"package/repository"
	"package/middleware"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	log.SetOutput(gin.DefaultWriter)

	// connect to the react rendering server
	goclient.Connect("tcp", "0.0.0.0:9000")

	defer repository.Connect(os.Getenv("DB_URL")).Close()
	repository.Migrate()

	middleware.InitAuthMiddleware()

	router := gin.Default()
	router.Static("/css", "./assets/css")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	controller.LoginController(router)
	controller.RegisterController(router)
	controller.UploadController(router)
	controller.HomeController(router)

	router.Run("0.0.0.0:8080")
}
