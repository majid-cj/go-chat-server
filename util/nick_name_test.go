package util

import (
	"fmt"
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateNickNameFromEmail(t *testing.T) {
	emails := []string{
		"test@email.com",
		"new@email.com",
		"mustmatch@email.com",
		"willtakefromhere@email.com",
	}

	for _, email := range emails {
		nickName := GenerateNickNameFromEmail(email)
		fmt.Printf("nickname %+v\n", nickName)
		assert.True(t, strings.ContainsAny(nickName, email))
	}
}
