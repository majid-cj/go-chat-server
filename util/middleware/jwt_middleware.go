package middleware

import (
	"github.com/majid-cj/go-chat-server/infrastructure/auth"
	"github.com/majid-cj/go-chat-server/util"

	"github.com/kataras/iris/v12"
)

// AuthenticationJWTMiddleware ...
func AuthenticationJWTMiddleware(c iris.Context) {
	err := auth.TokenValid(c.Request())
	if err != nil {
		util.ResponseError(util.GetError("unauthorized_access"), iris.StatusUnauthorized, c)
		return
	}
	c.Next()
}
