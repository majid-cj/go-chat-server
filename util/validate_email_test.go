package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateEmail(t *testing.T) {
	errNil := ValidateEmail("majid@mail.com")

	assert.Nil(t, errNil)

	err := ValidateEmail("majid.mail.com")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid format")
}
