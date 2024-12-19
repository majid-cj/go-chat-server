package routers

import (
	"github.com/majid-cj/go-chat-server/config"
	"github.com/majid-cj/go-chat-server/domain/entity"
	"github.com/majid-cj/go-chat-server/util"

	"github.com/kataras/iris/v12"
)

// MemberRouters ...
type MemberRouters struct {
	Config *config.AppConfig
}

// NewMemberRouters ...
func NewMemberRouters(config *config.AppConfig) *MemberRouters {
	return &MemberRouters{
		Config: config,
	}
}

// GetAllMembers ...
func (router *MemberRouters) GetAllMembers(c iris.Context) {
	var members entity.Members
	members, err := router.Config.Persistence.Member.GetMembers()
	if err != nil {
		util.ResponseError(err, iris.StatusNotFound, c)
		return
	}
	util.Response(members.GetMembersSerializer(), iris.StatusOK, c)
}

// GetMembersByType ...
func (router *MemberRouters) GetMembersByType(c iris.Context) {
	var members entity.Members
	members, err := router.Config.Persistence.Member.GetMembersByType(c.Params().GetUint8Default("type", 1))
	if err != nil {
		util.ResponseError(err, iris.StatusNotFound, c)
		return
	}
	util.Response(members.GetMembersSerializer(), iris.StatusOK, c)
}
