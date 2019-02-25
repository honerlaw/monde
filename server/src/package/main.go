package main

import (
	"github.com/gin-gonic/gin"
	"github.com/musawirali/preact-rpc/goclient"
	"package/controller"
	"github.com/joho/godotenv"
	"os"
	"package/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// connect to the react rendering server
	goclient.Connect("tcp", "0.0.0.0:9000")

	defer repository.Connect(os.Getenv("DB_URL")).Close()
	repository.Migrate()

	router := gin.Default()
	router.Static("/css", "../assets/css")
	router.StaticFile("/favicon.ico", "../assets/favicon.ico")

	controller.LoginController(router)
	controller.RegisterController(router)

	router.Run("0.0.0.0:8080")
}
