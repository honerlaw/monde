package user

import (
	"github.com/gin-gonic/gin"
	"server/user/route"
	"server/user/middleware"
	"server/user/model"
	"server/core/repository"
)

func RegisterRoutes(router *gin.Engine) {
	user := router.Group("/user")

	user.GET("/login", route.Login)
	user.GET("/register", route.Register)
	user.POST("/login", middleware.GetJWTAuth().LoginHandler)
	user.GET("/logout", route.Logout)
	user.POST("/register", route.RegisterPost)
}

func Migrate() {
	(&model.User{}).Migrate(repository.DB, repository.MigrateModel)
}
