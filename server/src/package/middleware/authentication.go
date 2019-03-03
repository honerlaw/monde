package middleware

import (
	"github.com/appleboy/gin-jwt"
	"time"
	"github.com/gin-gonic/gin"
	"os"
	"package/service"
	"errors"
	"net/http"
	"package/util"
)

type AuthPayload struct {
	ID    float64
	Roles []string
}

const identityKey = "ID"

var unauthorizedUrlToPageMap = gin.H{
	"/login":    "LoginPage",
	"/register": "RegisterPage",
}

var AuthMiddleware *jwt.GinJWTMiddleware

func GetAuthPayload(c *gin.Context) (*AuthPayload) {
	mw := AuthMiddleware

	claims, err := mw.GetClaimsFromJWT(c)
	if err != nil {
		return nil
	}

	if claims["exp"] == nil {
		return nil
	}

	if _, ok := claims["exp"].(float64); !ok {
		return nil
	}

	if int64(claims["exp"].(float64)) < mw.TimeFunc().Unix() {
		return nil
	}

	c.Set("JWT_PAYLOAD", claims)
	identity := mw.IdentityHandler(c)

	if identity != nil {
		c.Set(mw.IdentityKey, identity)
	}

	if !mw.Authorizator(identity, c) {
		return nil
	}

	return identity.(*AuthPayload)
}

func InitAuthMiddleware() {
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       os.Getenv("JWT_REALM"),
		Key:         []byte(os.Getenv("JWT_SECRET_KEY")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*AuthPayload); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					"Roles":     v.Roles,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			// convert the Roles to a string array
			intRoles := claims["Roles"].([]interface{})
			roles := make([]string, len(intRoles))
			for i, v := range intRoles {
				roles[i] = v.(string)
			}

			return &AuthPayload{
				ID:    claims[identityKey].(float64),
				Roles: roles,
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var req service.VerifyUserRequest

			if err := c.ShouldBind(&req); err != nil {
				return nil, errors.New("all fields are requred")
			}

			user, err := service.VerifyUser(req)
			if err != nil {
				return nil, err
			}

			return &AuthPayload{
				ID:    float64(user.ID),
				Roles: []string{"user"},
			}, nil
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			page := unauthorizedUrlToPageMap[c.Request.URL.String()]

			if page == nil {
				page = "UnauthorizedPage"
			}

			util.RenderPage(c, http.StatusUnauthorized, page.(string), gin.H{
				"error": message,
			})
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
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
	})

	if err != nil {
		panic(err)
	}

	AuthMiddleware = middleware
}
