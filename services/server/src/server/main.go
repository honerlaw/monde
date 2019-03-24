package main

import (
	"github.com/gin-gonic/gin"
	"github.com/musawirali/preact-rpc/goclient"
	"github.com/joho/godotenv"
	"log"
	"server/service/aws"
	"server/middleware/auth"
	"server/repository"
	"server/routes/user"
	"server/routes"
	"net/http"
	"os"
	"server/media"
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
	repository.Migrate()

	router := gin.Default()

	// initialize the auth middleware
	auth.Init(router)

	router.Static("/css/", "./assets/css/")
	router.Static("/js/", "./assets/js/")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	media.RegisterRoutes(router)
	router.GET("/user/login", user.Login)
	router.GET("/user/register", user.Register)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})
	router.GET("/", routes.Home)

	router.Run("0.0.0.0:" + os.Getenv("PORT"))
}
