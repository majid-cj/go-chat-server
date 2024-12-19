package util

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	nickNameRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)
)

// GenerateNickNameFromEmail ...
func GenerateNickNameFromEmail(email string) string {
	ranULID := ULID()
	nickName := nickNameRegex.ReplaceAllString(strings.Split(email, "@")[0], "")
	trimLength := len(nickName)

	if trimLength > 8 {
		trimLength /= 2
		nickName = nickName[:trimLength]
	}
	return fmt.Sprintf("%s%s", nickName, ranULID[:trimLength])
}
