package entity

import (
	"errors"
	"regexp"

	"github.com/majid-cj/go-chat-server/util"
	"github.com/majid-cj/go-chat-server/util/security"
)

// SignUp ...
type SignUp struct {
	DisplayName  string `json:"display_name"`
	UniqueId     string `json:"unique_id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
}

// ValidateSignUpMember ...
func (signUp *SignUp) ValidateSignUpMember() error {
	if err := util.ValidateEmail(signUp.Email); err != nil {
		return errors.New("invalid_email")
	}
	if _, err := regexp.Match(security.PASSWORD_PATTERN, []byte(signUp.Password)); err != nil {
		return errors.New("invalid password")
	}
	return nil
}

// ValidateSocialMember ...
func (signUp *SignUp) ValidateSocialMember() error {
	if err := util.ValidateEmail(signUp.Email); err != nil {
		return errors.New("invalid_email")
	}
	if signUp.Password != "" {
		return errors.New("invalid_password")
	}
	return nil
}
