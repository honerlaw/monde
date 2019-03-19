package auth

import (
	"github.com/appleboy/gin-jwt"
	"time"
	"github.com/gin-gonic/gin"
	"os"
	"errors"
	"net/http"
	"package/util"
	"package/service/user"
)

// @todo split this up / re-organize it?

type AuthPayload struct {
	ID    int64
	Roles []string
}

const identityKey = "ID"

var unauthorizedUrlToPageMap = gin.H{
	"/user/login":    "LoginPage",
	"/user/register": "RegisterPage",
}

func createJwtMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
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
				ID:    int64(claims[identityKey].(float64)),
				Roles: roles,
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var req user.VerifyRequest

			if err := c.ShouldBind(&req); err != nil {
				return nil, errors.New("all fields are requred")
			}

			verifiedUser, err := user.Verify(req)
			if err != nil {
				return nil, err
			}

			return &AuthPayload{
				ID:    int64(verifiedUser.ID), // definition says unit, runtime says float64...
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
}

func logoutHandler(mw *jwt.GinJWTMiddleware, c *gin.Context) {
	cookieName := mw.CookieName
	cookie, err := c.Request.Cookie(cookieName)

	if err != nil {
		panic(err)
	}

	cookie.Value = "invalid"
	cookie.Expires = time.Unix(0, 0)

	http.SetCookie(c.Writer, cookie)

	c.Redirect(http.StatusFound, "/")
	c.Abort()
}

func handleRegister(mw *jwt.GinJWTMiddleware, c *gin.Context) {
	var req user.CreateRequest

	if err := c.ShouldBind(&req); err != nil {
		util.RenderPage(c, http.StatusBadRequest, "RegisterPage", gin.H{
			"error": "all fields are required",
		})
		return
	}

	_, err := user.Create(req)

	if err != nil {
		util.RenderPage(c, http.StatusBadRequest, "RegisterPage", gin.H{
			"error": err.Error(),
		})
		return
	}

	mw.LoginHandler(c)
}

func Init(router *gin.Engine) {
	mw, err := createJwtMiddleware()

	router.POST("/user/login", mw.LoginHandler)
	router.GET("/user/logout", func(c *gin.Context) {
		logoutHandler(mw, c)
	})
	router.POST("/user/register", func(c *gin.Context) {
		handleRegister(mw, c)
	})

	if err != nil {
		panic(err)
	}

	// we use custom middleware so we don't trigger the unauthorized callback, instead we will handle it ourselves
	router.Use(func(c *gin.Context) {
		// we always set this so each handler can check if it exists, whether or not are actually logged in
		c.Set("JWT_IDENTITY", nil)

		claims, err := mw.GetClaimsFromJWT(c)
		if err != nil {
			return
		}

		if claims["exp"] == nil {
			return
		}

		if _, ok := claims["exp"].(float64); !ok {
			return
		}

		if int64(claims["exp"].(float64)) < mw.TimeFunc().Unix() {
			return
		}

		c.Set("JWT_PAYLOAD", claims)
		identity := mw.IdentityHandler(c)

		if identity != nil {
			c.Set(mw.IdentityKey, identity)
		}

		if !mw.Authorizator(identity, c) {
			return
		}

		c.Set("JWT_IDENTITY", identity.(*AuthPayload))

		c.Next()
	})


}
