package util

import (
	"strings"
)

// VerifyCode ...
func VerifyCode() string {
	return strings.ToUpper(ULID()[:6])
}
