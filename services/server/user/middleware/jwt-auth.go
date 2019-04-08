package middleware

import (
	"sync"
	"services/server/user/service"
	"github.com/appleboy/gin-jwt"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"errors"
	"net/http"
	"services/server/core/render"
	"services/server/core/util"
)

type AuthPayload struct {
	ID    string
	Roles []string
}

var jwtAuthSync sync.Once
var jwtAuth *jwt.GinJWTMiddleware

func createJwtMiddleware() (*jwt.GinJWTMiddleware) {
	return &jwt.GinJWTMiddleware{
		Realm:      os.Getenv("JWT_REALM"),
		Key:        []byte(os.Getenv("JWT_SECRET_KEY")),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*AuthPayload); ok {
				return jwt.MapClaims{
					"id":    v.ID,
					"roles": v.Roles,
				}
			}
			return nil
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var req service.VerifyRequest

			if err := c.ShouldBind(&req); err != nil {
				return nil, errors.New("all fields are requred")
			}

			verifiedUser, err := service.Verify(req)
			if err != nil {
				return nil, err
			}

			return &AuthPayload{
				ID:    verifiedUser.ID, // definition says unit, runtime says float64...
				Roles: []string{"user"},
			}, nil
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			util.Redirect(c, "/")
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			var req service.VerifyRequest
			c.ShouldBind(&req);
			render.RenderPage(c, http.StatusUnauthorized, gin.H{
				"username": req.Username,
				"error":    message,
			})
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {

			// @todo this should check a lot more than it does

			if v, ok := data.(*AuthPayload); ok {
				for _, b := range v.Roles {
					if b == "user" {
						return true
					}
				}
			}
			return false
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		SendCookie:    true,
	}
}

func GetJWTAuth() (*jwt.GinJWTMiddleware) {
	jwtAuthSync.Do(func() {
		jwtAuth = createJwtMiddleware()
	})
	return jwtAuth
}
