package routers

import (
	"github.com/majid-cj/go-chat-server/config"
	"github.com/majid-cj/go-chat-server/domain/entity"
	"github.com/majid-cj/go-chat-server/util"

	"github.com/kataras/iris/v12"
)

// VerifyCodeRouter ...
type VerifyCodeRouter struct {
	Config *config.AppConfig
}

// NewVerifyCodeRouter ...
func NewVerifyCodeRouter(config *config.AppConfig) *VerifyCodeRouter {
	return &VerifyCodeRouter{
		Config: config,
	}
}

// NewVerifyCode ...
func (router *VerifyCodeRouter) NewVerifyCode(c iris.Context) {
	var verifyCode entity.VerificationCode
	err := c.ReadJSON(&verifyCode)
	if err != nil {
		util.ResponseError(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}
	verifyCode.PrepareVerificationCode(verifyCode.Member, verifyCode.CodeType)
	_, err = router.Config.Persistence.VerifyCode.CreateVerificationCode(&verifyCode)
	if err != nil {
		util.ResponseError(err, iris.StatusBadRequest, c)
		return
	}

	err = util.SendMail([]string{verifyCode.Email}, util.ReceiverMail{
		ReceiverMail: verifyCode.Email,
		ReceiverName: "",
		ReceiverCode: verifyCode.Code,
	}, "verification_code.txt", "New Verification Code")

	if err != nil {
		util.ResponseError(util.GetError("general_error"), iris.StatusBadRequest, c)
		return
	}
	util.Response(nil, iris.StatusCreated, c)
}

// VerificationCodeFromEmail ...
func (router *VerifyCodeRouter) VerificationCodeFromEmail(c iris.Context) {
	var verifyCode entity.VerificationCode
	err := c.ReadJSON(&verifyCode)
	if err != nil {
		util.ResponseError(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}
	_, err = router.Config.Persistence.VerifyCode.CreateVerificationCodeFromEmail(&verifyCode)
	if err != nil {
		util.ResponseError(err, iris.StatusBadRequest, c)
		return
	}

	err = util.SendMail([]string{verifyCode.Email}, util.ReceiverMail{
		ReceiverMail: verifyCode.Email,
		ReceiverName: "",
		ReceiverCode: verifyCode.Code,
	}, "verification_code.txt", "New Verification Code")

	if err != nil {
		util.ResponseError(util.GetError("general_error"), iris.StatusBadRequest, c)
		return
	}

	util.Response(nil, iris.StatusCreated, c)
}

// ResetPasswordVerifyCode ...
func (router *VerifyCodeRouter) ResetPasswordVerifyCode(c iris.Context) {
	var verifyCode entity.VerificationCode
	err := c.ReadJSON(&verifyCode)
	if err != nil {
		util.ResponseError(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}

	err = router.Config.Persistence.VerifyCode.ResetPassword(&verifyCode)
	if err != nil {
		util.ResponseError(err, iris.StatusBadRequest, c)
		return
	}
	util.Response(nil, iris.StatusOK, c)
}

// CheckVerifyCode ...
func (router *VerifyCodeRouter) CheckVerifyCode(c iris.Context) {
	var verifyCode entity.VerificationCode
	err := c.ReadJSON(&verifyCode)
	if err != nil {
		util.ResponseError(util.GetError("error_parsing_data"), iris.StatusUnprocessableEntity, c)
		return
	}

	err = router.Config.Persistence.VerifyCode.CheckVerificationCode(&verifyCode)
	if err != nil {
		util.ResponseError(err, iris.StatusBadRequest, c)
		return
	}
	util.Response(nil, iris.StatusOK, c)
}

// RenewVerifyCode ...
func (router *VerifyCodeRouter) RenewVerifyCode(c iris.Context) {
	var verifyCode entity.VerificationCode
	err := c.ReadJSON(&verifyCode)

	if err != nil {
		util.ResponseError(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}
	_, err = router.Config.Persistence.VerifyCode.RenewVerificationCode(&verifyCode)
	if err != nil {
		util.ResponseError(err, iris.StatusBadRequest, c)
		return
	}

	err = util.SendMail([]string{verifyCode.Email}, util.ReceiverMail{
		ReceiverMail: verifyCode.Email,
		ReceiverName: "",
		ReceiverCode: verifyCode.Code,
	}, "verification_code.txt", "Renew Verification Code")

	if err != nil {
		util.ResponseError(util.GetError("general_error"), iris.StatusBadRequest, c)
		return
	}

	util.Response(nil, iris.StatusCreated, c)
}
