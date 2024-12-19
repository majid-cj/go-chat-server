package routers

import (
	"os"
	"time"

	"github.com/majid-cj/go-chat-server/config"
	"github.com/majid-cj/go-chat-server/domain/entity"
	"github.com/majid-cj/go-chat-server/infrastructure/auth"
	"github.com/majid-cj/go-chat-server/util"
	"github.com/majid-cj/go-chat-server/util/security"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12"
)

// AuthenticationRouter ...
type AuthenticationRouter struct {
	Config *config.AppConfig
}

// NewAuthenticationRouter ...
func NewAuthenticationRouter(config *config.AppConfig) *AuthenticationRouter {
	return &AuthenticationRouter{
		Config: config,
	}
}

// SignUp ...
func (router *AuthenticationRouter) SignUp(c iris.Context) {
	var member entity.Member
	var signUp entity.SignUp
	var profile entity.MemberProfile
	var code entity.VerificationCode

	err := c.ReadJSON(&signUp)
	if err != nil {
		util.ResponseError(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}

	member.Email = signUp.Email
	err = signUp.ValidateSignUpMember()
	if err != nil {
		util.ResponseError(err, iris.StatusUnprocessableEntity, c)
		return
	}
	member.PrepareMember()
	member.Password = security.NewPassHash(member.ID, signUp.Email, signUp.Password, []byte{128})
	newMember, err := router.Config.Persistence.Member.CreateMember(&member)
	if err != nil {
		util.ResponseError(err, iris.StatusBadRequest, c)
		return
	}
	profile.PrepareDefaultProfile(
		newMember.ID,
		signUp.DisplayName,
		member.Email)
	newProfile, err := router.Config.Persistence.Profile.CreateMemberProfile(&profile)

	if err != nil || newProfile == nil {
		if err := router.Config.Persistence.Member.DeleteMember(member.ID); err != nil {
			util.ResponseError(err, iris.StatusBadRequest, c)
			return
		}
	}

	token, err := router.Config.Token.CreateJWTToken(
		newMember.ID,
		newProfile.ID,
		signUp.UniqueId,
	)
	if err != nil {
		util.ResponseError(util.GetError("general_error"), iris.StatusUnprocessableEntity, c)
		return
	}

	_, createError := router.Config.Auth.Auth.CreateToken(newMember.ID, token)
	if createError != nil {
		util.ResponseError(util.GetError("general_error"), iris.StatusBadRequest, c)
		return
	}

	code.PrepareVerificationCode(member.ID, 1)
	_, err = router.Config.Persistence.VerifyCode.CreateVerificationCode(&code)
	if err != nil {
		util.ResponseError(util.GetError("general_error"), iris.StatusBadRequest, c)
		return
	}

	err = util.SendMail([]string{newMember.GetMemberSerializer().Email}, util.ReceiverMail{
		ReceiverMail: newMember.GetMemberSerializer().Email,
		ReceiverName: newProfile.DisplayName,
		ReceiverCode: code.Code,
	}, "verification_code.txt", "Sign Up Verification Code")

	if err != nil {
		util.ResponseError(util.GetError("general_error"), iris.StatusBadRequest, c)
		return
	}

	memberResponse := make(map[string]interface{})

	memberResponse["token"] = token
	memberResponse["member"] = newMember.GetMemberSerializer()
	memberResponse["profile"] = newProfile

	util.Response(memberResponse, iris.StatusCreated, c)
}

// SignIn ...
func (router *AuthenticationRouter) SignIn(c iris.Context) {
	var member *entity.SignUp
	var code entity.VerificationCode

	err := c.ReadJSON(&member)
	if err != nil {
		util.ResponseError(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}
	err = member.ValidateSignUpMember()
	if err != nil {
		util.ResponseError(err, iris.StatusUnprocessableEntity, c)
		return
	}
	memberLogin, err := router.Config.Persistence.Member.GetMemberByEmailAndPassword(member)
	if err != nil {
		util.ResponseError(err, iris.StatusNotFound, c)
		return
	}
	profile, err := router.Config.Persistence.Profile.GetMemberProfileByMemberID(memberLogin.ID)
	if err != nil {
		util.ResponseError(err, iris.StatusNotFound, c)
		return
	}
	token, err := router.Config.Token.CreateJWTToken(
		memberLogin.ID,
		profile.ID,
		member.UniqueId,
	)
	if err != nil {
		util.ResponseError(err, iris.StatusBadRequest, c)
		return
	}

	loginSession, err := router.Config.Auth.Auth.CreateToken(memberLogin.ID, token)
	if err != nil {
		util.ResponseError(err, iris.StatusBadRequest, c)
		return
	}

	if len(loginSession) > 0 {
		ip, _ := router.Config.IPInfo.GetIPAddr()
		_ = util.SendMail([]string{memberLogin.GetMemberSerializer().Email}, util.ActiveLoginMail{
			ReceiverMail: memberLogin.GetMemberSerializer().Email,
			ReceiverName: profile.DisplayName,
			LoginTime:    time.Now().UTC().Format(time.ANSIC),
			DeviceInfo:   c.Request().Header.Get("User-Agent"),
			IPAddress:    ip,
		}, "active_login.txt", "Sign in from different device")
	}

	if !memberLogin.Verified {
		code.PrepareVerificationCode(memberLogin.ID, 1)
		_, err = router.Config.Persistence.VerifyCode.CreateVerificationCode(&code)
		if err != nil {
			util.ResponseError(util.GetError("general_error"), iris.StatusBadRequest, c)
			return
		}

		_ = util.SendMail([]string{memberLogin.GetMemberSerializer().Email}, util.ReceiverMail{
			ReceiverMail: memberLogin.GetMemberSerializer().Email,
			ReceiverName: profile.DisplayName,
			ReceiverCode: code.Code,
		}, "verification_code.txt", "New Login Verification Code")

		if err != nil {
			util.ResponseError(util.GetError("general_error"), iris.StatusBadRequest, c)
			return
		}
	}

	memberResponse := make(map[string]interface{})
	memberResponse["token"] = token
	memberResponse["member"] = memberLogin.GetMemberSerializer()
	memberResponse["profile"] = profile

	util.Response(memberResponse, iris.StatusOK, c)
}

// Logout ...
func (router *AuthenticationRouter) Logout(c iris.Context) {
	token, err := router.Config.Token.ExtractJWTTokenMetadata(c.Request(), true)
	if err != nil {
		util.ResponseError(err, iris.StatusUnauthorized, c)
		return
	}
	// router.Config.Auth.Auth.DeleteAccessToken(token)
	err = router.Config.Auth.Auth.DeleteAccessToken(token)
	if err != nil {
		util.ResponseError(err, iris.StatusBadRequest, c)
		return
	}
	util.Response(nil, iris.StatusOK, c)
}

// Refresh ...
func (router *AuthenticationRouter) Refresh(c iris.Context) {
	var data struct {
		Refresh string `json:"refresh"`
	}

	err := c.ReadJSON(&data)
	if err != nil {
		util.ResponseError(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}

	refreshToken := data.Refresh
	token, _ := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, util.GetError("general_error")
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if _, ok := token.Claims.(jwt.Claims); !ok {
		util.ResponseError(err, iris.StatusUnauthorized, c)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		userId, ok := claims["user_id"].(string)
		if !ok {
			util.ResponseError(util.GetError("general_error"), iris.StatusUnauthorized, c)
			return
		}

		profileId, ok := claims["profile_id"].(string)
		if !ok {
			util.ResponseError(util.GetError("general_error"), iris.StatusUnauthorized, c)
			return
		}

		uniqueId, ok := claims["unique_id"].(string)
		if !ok {
			util.ResponseError(util.GetError("general_error"), iris.StatusUnauthorized, c)
			return
		}

		refreshUUID, ok := claims["refresh_uuid"].(string)
		if !ok {
			util.ResponseError(util.GetError("general_error"), iris.StatusUnauthorized, c)
			return
		}
		token, err := router.Config.Token.ExtractJWTTokenMetadata(c.Request(), false)
		if err != nil {
			util.ResponseError(err, iris.StatusUnauthorized, c)
			return
		}
		router.Config.Auth.Auth.DeleteAccessToken(token)
		router.Config.Auth.Auth.DeleteRefreshToken(refreshUUID)
		newToken, err := router.Config.Token.CreateJWTToken(userId, profileId, uniqueId)
		if err != nil {
			util.ResponseError(err, iris.StatusUnauthorized, c)
			return
		}

		_, err = router.Config.Auth.Auth.CreateToken(userId, newToken)
		if err != nil {
			util.ResponseError(err, iris.StatusUnauthorized, c)
			return
		}

		util.Response(newToken, iris.StatusOK, c)
	} else {
		util.ResponseError(util.GetError("general_error"), iris.StatusUnauthorized, c)
		return
	}
}

// UpdatePassword ...
func (router *AuthenticationRouter) UpdatePassword(c iris.Context) {
	var data struct {
		Password        string `json:"password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	err := c.ReadJSON(&data)
	if err != nil {
		util.ResponseError(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}

	memberId := auth.ExtractTokenClaims(c.Request(), "user_id")
	member, err := router.Config.Persistence.Member.GetMember(memberId)
	if err != nil {
		util.ResponseError(err, iris.StatusNotFound, c)
		return
	}
	valid := security.EqualPassHash(member.ID, member.Email, data.Password, member.Password)
	if !valid {
		util.ResponseError(err, iris.StatusUnauthorized, c)
		return
	}

	if data.NewPassword != data.ConfirmPassword {
		util.ResponseError(util.GetError("password_mismatch"), iris.StatusBadRequest, c)
		return
	}

	member.Password = security.NewPassHash(member.ID, member.Email, data.ConfirmPassword, []byte{128})
	member.UpdateAt = util.GetTimeNow()

	err = router.Config.Persistence.Member.UpdatePassword(member)
	if err != nil {
		util.ResponseError(err, iris.StatusBadRequest, c)
		return
	}

	util.Response(true, iris.StatusOK, c)
}
