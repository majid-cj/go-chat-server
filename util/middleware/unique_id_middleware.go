package middleware

import (
	"github.com/majid-cj/go-chat-server/infrastructure/auth"
	"github.com/majid-cj/go-chat-server/util"

	"github.com/kataras/iris/v12"
)

// UniqueIdMiddleware ...
func UniqueIdMiddleware(c iris.Context) {
	uniqueId := auth.ExtractTokenClaims(c.Request(), "unique_id")
	requestUniqueId := c.Request().Header.Get("UniqueId")
	if uniqueId != requestUniqueId {
		util.ResponseError(util.GetError("unauthorized_access"), iris.StatusUnauthorized, c)
		return
	}
	c.Next()
}
