package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"server/media"
	"server/core/repository"
	"server/core/service/aws"
	"server/user"
	"server/core"
	"server/user/middleware"
	mediaMW "server/media/middleware"
	renderMW "server/core/render/middleware"
	"server/core/render"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(gin.DefaultWriter)

	// init aws session
	err = aws.InitSession()
	if err != nil {
		log.Fatal(err)
	}

	defer repository.GetRepository().Migrate().DB().Close()

	router := gin.Default()

	router.Use(middleware.GetJWTAuth().MiddlewareFunc())
	router.Use(renderMW.ReactRenderMiddleware("./assets/js/bundle.js", false, router))
	router.Use(mediaMW.UploadFormMiddleware())

	router.Static("/css/", "./assets/css/")
	router.Static("/js/", "./assets/js/")
	router.Static("/gen/", "./assets/gen/")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	router.NoRoute(render.RenderNoRoute)

	user.RegisterRoutes(router)
	core.RegisterRoutes(router)
	media.RegisterRoutes(router)

	router.Run("0.0.0.0:" + os.Getenv("PORT"))
}