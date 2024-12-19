package util

import (
	"strings"
)

// EscapeString ...
func EscapeString(str string) string {
	return strings.TrimSpace(strings.ToLower(str))
}
