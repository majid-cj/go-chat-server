package router

import (
	"github.com/majid-cj/go-chat-server/config"
	"github.com/majid-cj/go-chat-server/router/routers"
	"github.com/majid-cj/go-chat-server/util/middleware"
)

// APIVersionOne ...
func APIVersionOne(appConfig *config.AppConfig) {
	authentication := routers.NewAuthenticationRouter(appConfig)
	member := routers.NewMemberRouters(appConfig)
	verifyCode := routers.NewVerifyCodeRouter(appConfig)
	profile := routers.NewMemberProfileRouter(appConfig)
	chat := routers.NewChatRouter(appConfig)

	appConfig.App.UseGlobal(middleware.RateLimit)

	apiV1 := appConfig.App.Party("/api/v1")
	{
		apiV1.Get("/{nick_name:string}", middleware.AuthenticationJWTMiddleware, middleware.UniqueIdMiddleware, profile.GetProfileByNickName)

		apiV1.Get("/chat-list", middleware.AuthenticationJWTMiddleware, middleware.UniqueIdMiddleware, chat.GetChatList)
		apiV1.Get("/chat-counter", middleware.AuthenticationJWTMiddleware, middleware.UniqueIdMiddleware, chat.GetChatCounter)

		apiV1.Get("/ws/{sender:string}/{receiver:string}", chat.HandleRequest)
		appConfig.Melody.HandleConnect(chat.HandleConnect)
		appConfig.Melody.HandleMessage(chat.HandleMessage)
		appConfig.Melody.HandleDisconnect(chat.HandleDisconnect)
		appConfig.Melody.HandleClose(chat.HandleClose)

		MemberRouteEndPoints(authentication, member, verifyCode, apiV1)

	}
}
