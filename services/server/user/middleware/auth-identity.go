package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthIdentity() (gin.HandlerFunc) {
	return func(c *gin.Context) {
		mw := GetJWTAuth()

		c.Set("JWT_AUTH_PAYLOAD", nil)

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

		c.Set("JWT_AUTH_PAYLOAD", identity.(*AuthPayload))

		c.Next()
	}
}
