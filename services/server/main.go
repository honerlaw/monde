package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"services/server/media"
	"services/server/core/repository"
	"services/server/core/service/aws"
	"services/server/user"
	"services/server/core"
	mediaMW "services/server/media/middleware"
	renderMW "services/server/core/render/middleware"
	"services/server/core/render"
	"services/server/user/middleware"
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

	// initialize each module
	userModule := user.Init()
	mediaModule := media.Init(userModule.ChannelService)

	// initialize the jwt auth middleware
	middleware.InitJWTAuth(userModule.UserService)

	// start setting up routes
	router := gin.Default()

	router.Use(middleware.AuthIdentity())
	router.Use(renderMW.ReactRenderMiddleware("./assets/js/bundle.js", os.Getenv("REACT_POOL_TYPE") == "on_demand", router))
	router.Use(mediaMW.UploadFormMiddleware(userModule.ChannelService))

	router.Static("/css/", "./assets/css/")
	router.Static("/js/", "./assets/js/")
	router.Static("/gen/", "./assets/gen/")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	router.NoRoute(render.RenderNoRoute)

	// register the module specific routes last
	userModule.RegisterRoutes(router)
	core.RegisterRoutes(router)
	mediaModule.RegisterRoutes(router)

	router.Run("0.0.0.0:" + os.Getenv("PORT"))
}
