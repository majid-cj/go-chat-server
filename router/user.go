package router

import (
	"github.com/kataras/iris/v12/core/router"
	"github.com/majid-cj/go-chat-server/router/routers"
	"github.com/majid-cj/go-chat-server/util/middleware"
)

// MemberRouteEndPoints ...
func MemberRouteEndPoints(
	authentication *routers.AuthenticationRouter,
	member *routers.MemberRouters,
	verifyCode *routers.VerifyCodeRouter,
	APIVersion router.Party,
) {
	userRoute := APIVersion.Party("/user")
	{
		userRoute.Post("/sign-up", authentication.SignUp)
		userRoute.Post("/sign-in", authentication.SignIn)

		userRoute.Post("/reset/code", verifyCode.VerificationCodeFromEmail)
		userRoute.Post("/reset/password", verifyCode.ResetPasswordVerifyCode)

		userRoute.Post("/refresh", authentication.Refresh)

		userRoute.Use(middleware.AuthenticationJWTMiddleware, middleware.UniqueIdMiddleware)
		userRoute.Put("/password", authentication.UpdatePassword)
		userRoute.Post("/logout", authentication.Logout)

		userRoute.Post("/verify/code", verifyCode.NewVerifyCode)
		userRoute.Post("/verify/check", verifyCode.CheckVerifyCode)
		userRoute.Post("/verify/renew", verifyCode.RenewVerifyCode)
		userRoute.Post("/verify/reset/password", verifyCode.ResetPasswordVerifyCode)

	}
}
