package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_EqualPassHash ...
func Test_EqualPassHash(t *testing.T) {
	salt := NewPassHash("random-id-1234", "test@mail.com", "123456", []byte{128})

	assert.True(t, EqualPassHash("random-id-1234", "test@mail.com", "123456", salt))

	assert.False(t, EqualPassHash("random-id-1234", "test@mail.com", "123455", salt))
}
