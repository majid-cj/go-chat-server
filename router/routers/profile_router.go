package routers

import (
	"strings"

	"github.com/majid-cj/go-chat-server/config"
	"github.com/majid-cj/go-chat-server/domain/entity"
	"github.com/majid-cj/go-chat-server/infrastructure/auth"
	"github.com/majid-cj/go-chat-server/util"
	"github.com/samber/lo"

	"github.com/albrow/forms"
	"github.com/kataras/iris/v12"
)

// MemberProfileRouter ...
type MemberProfileRouter struct {
	Config *config.AppConfig
}

// NewMemberProfileRouter ...
func NewMemberProfileRouter(config *config.AppConfig) *MemberProfileRouter {
	return &MemberProfileRouter{
		Config: config,
	}
}

// UpdateMemberProfile ...
func (router *MemberProfileRouter) UpdateMemberProfile(c iris.Context) {
	var profile entity.MemberProfile
	var fileURL string

	data, err := forms.Parse(c.Request())
	if err != nil {
		util.ResponseError(util.GetError("error_parsing_data"), iris.StatusUnprocessableEntity, c)
		return
	}

	validator := data.Validator()
	validator.Require("display_name")
	validator.Require("nick_name")
	if validator.HasErrors() {
		util.ResponseError(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}

	displayName := data.Get("display_name")
	nickName := data.Get("nick_name")

	profile.NickName = strings.ToLower(nickName)
	profile.DisplayName = displayName

	err = profile.ValidateMemberProfile()
	if err != nil {
		util.ResponseError(err, iris.StatusBadRequest, c)
		return
	}

	updateProfile, err := router.Config.Persistence.Profile.GetMemberProfileByID(c.Params().GetString("id"))
	if err != nil {
		util.ResponseError(util.GetError("general_error"), iris.StatusBadRequest, c)
		return
	}

	file, fileHeader, err := c.FormFile("profile_image")
	if err == nil {
		fileURL, err = router.Config.Upload.UploadFile(fileHeader, file, "profile")
		if err != nil {
			util.ResponseError(err, iris.StatusBadRequest, c)
			return
		}
	} else {
		fileURL = updateProfile.ProfileImage
	}

	updateProfile.ProfileImage = fileURL
	updateProfile.NickName = nickName
	updateProfile.DisplayName = displayName
	updatedProfile, err := router.Config.Persistence.Profile.UpdateMemberProfile(updateProfile)

	if err != nil {
		util.ResponseError(err, iris.StatusBadRequest, c)
		return
	}

	util.Response(updatedProfile, iris.StatusOK, c)
}

// GetProfileByMember ...
func (router *MemberProfileRouter) GetProfileByMember(c iris.Context) {
	member := auth.ExtractTokenClaims(c.Request(), "user_id")
	profile, err := router.Config.Persistence.Profile.GetMemberProfileByMemberID(member)
	if err != nil {
		util.ResponseError(err, iris.StatusNotFound, c)
		return
	}
	util.Response(profile, iris.StatusOK, c)
}

// GetProfile ...
func (router *MemberProfileRouter) GetProfile(c iris.Context) {
	profileId := auth.ExtractTokenClaims(c.Request(), "profile_id")
	profile, err := router.Config.Persistence.Profile.GetMemberProfileByID(profileId)
	if err != nil {
		util.ResponseError(err, iris.StatusNotFound, c)
		return
	}
	util.Response(profile, iris.StatusOK, c)
}

// GetProfileByID ...
func (router *MemberProfileRouter) GetProfileByID(c iris.Context) {
	profileId := c.Params().Get("profile")
	profile, err := router.Config.Persistence.Profile.GetMemberProfileByID(profileId)
	if err != nil {
		util.ResponseError(err, iris.StatusNotFound, c)
		return
	}
	util.Response(profile, iris.StatusOK, c)
}

// GetProfileByNickName ...
func (router *MemberProfileRouter) GetProfileByNickName(c iris.Context) {
	if !lo.Contains([]string{"search", "qr_code"}, c.URLParam("source")) {
		util.ResponseError(util.GetError("unauthorized_access"), iris.StatusUnauthorized, c)
		return
	}

	profile, err := router.Config.Persistence.Profile.GetMemberProfileByNickName(strings.ToLower(c.Params().Get("nick_name")))
	if err != nil {
		util.ResponseError(err, iris.StatusNotFound, c)
		return
	}
	util.Response(profile, iris.StatusOK, c)
}
