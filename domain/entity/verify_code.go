package entity

import (
	"time"

	"github.com/majid-cj/go-chat-server/util"
)

// VerifyCodeType ...
var VerifyCodeType = map[uint8]string{
	1: "New User",
	2: "Renew Password",
	3: "Reset Password",
}

// VerificationCode ...
type VerificationCode struct {
	ID        string    `bson:"id" json:"id"`
	Member    string    `bson:"member" json:"member"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"password"`
	Code      string    `bson:"code" json:"code"`
	Taken     bool      `bson:"taken" json:"taken"`
	CodeType  uint8     `bson:"code_type" json:"code_type"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	ExpiredAt time.Time `bson:"expired_at" json:"expired_at"`
}

// PrepareVerificationCode ...
func (vc *VerificationCode) PrepareVerificationCode(member string, codeType uint8) {
	vc.ID = util.ULID()
	vc.Member = member
	vc.Code = util.VerifyCode()
	vc.Taken = false
	vc.CodeType = codeType
	vc.CreatedAt = util.GetTimeNow()
	vc.ExpiredAt = util.TimeAfter(vc.CreatedAt, (time.Minute * 5))
}
