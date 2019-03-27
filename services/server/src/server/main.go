package main

import (
	"github.com/gin-gonic/gin"
	"github.com/musawirali/preact-rpc/goclient"
	"github.com/joho/godotenv"
	"log"
	"os"
	"server/media"
	"server/core/repository"
	"server/core/service/aws"
	"server/user"
	"server/core"
	"server/user/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	log.SetOutput(gin.DefaultWriter)

	// init aws session
	err = aws.InitSession()
	if err != nil {
		panic(err)
	}

	// connect to the react rendering server
	goclient.Connect("tcp", "0.0.0.0:9000")

	defer repository.Connect().Close()

	user.Migrate()
	media.Migrate()

	router := gin.Default()

	router.Use(middleware.AuthIdentity())

	router.Static("/css/", "./assets/css/")
	router.Static("/js/", "./assets/js/")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	media.RegisterRoutes(router)
	user.RegisterRoutes(router)
	core.RegisterRoutes(router)

	router.Run("0.0.0.0:" + os.Getenv("PORT"))
}
