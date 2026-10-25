package transport

import (
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

type SessionCustomClaims struct {
	Role string `json:"role"`
}

func SessionClaimMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
		if !ok {
			Fail(c, http.StatusUnauthorized, "unauthorized request")
			c.Abort()
		}
		// this part is currently a placeholder, will implement role check and tier
		// check in the future
		_, ok = claims.Custom.(*SessionCustomClaims)
		if !ok {
			Fail(c, http.StatusUnauthorized, "failed to validate user info")
			c.Abort()
		}
		c.Next()
	}
}
