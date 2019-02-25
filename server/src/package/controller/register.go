package controller

import (
	"github.com/gin-gonic/gin"
	"package/util"
	"net/http"
	"package/repository"
	"package/model"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username       string `form:"username" binding:"required"`
	Password       string `form:"password" binding:"required"`
	VerifyPassword string `form:"verify_password" binding:"required"`
}

func RegisterController(router *gin.Engine) {
	router.GET("/register", renderRegisterPage)
	router.POST("/register", handleRegister)
}

func renderRegisterPage(c *gin.Context) {
	util.RenderPage(c, http.StatusOK, "RegisterPage", nil)
}

func handleRegister(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		util.RenderPage(c, http.StatusBadRequest, "RegisterPage", gin.H{
			"error": "all fields are required",
		})
		return
	}

	if len(req.Username) < 6 {
		util.RenderPage(c, http.StatusBadRequest, "RegisterPage", gin.H{
			"error": "username must be at least 6 characters in length",
		})
		return
	}

	if len(req.Password) < 6 {
		util.RenderPage(c, http.StatusBadRequest, "RegisterPage", gin.H{
			"error": "password must be at least 6 characters in length",
		})
		return
	}

	if req.Password != req.VerifyPassword {
		util.RenderPage(c, http.StatusBadRequest, "RegisterPage", gin.H{
			"error": "passwords do not match",
		})
		return
	}

	var user model.User
	repository.DB.Where(model.User{Username: req.Username}).First(&user)

	// no user was found
	if (model.User{}) != user {
		util.RenderPage(c, http.StatusBadRequest, "RegisterPage", gin.H{
			"error": "user already exists",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		util.RenderPage(c, http.StatusInternalServerError, "RegisterPage", gin.H{
			"error": "something went wrong. please try again",
		})
		return
	}

	user = model.User{
		Username: req.Username,
		Hash:     string(hash),
	}

	repository.DB.Save(&user);

	util.RenderPage(c, http.StatusOK, "RegisterPage", nil)
}
