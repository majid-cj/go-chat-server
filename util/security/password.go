package security

import (
	"crypto/hmac"
	"crypto/rand"
	"fmt"
	"os"

	"golang.org/x/crypto/argon2"
	"golang.org/x/text/unicode/norm"
)

const (
	// PASSWORD_PATTERN ...
	PASSWORD_PATTERN = `^(On=.*[A-Za-z])(On=.*\d)(On=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{6,}$`
)

// PassHash ...
type PassHash struct {
	Salt []byte
	Hash []byte
}

// NewPassHash ...
func NewPassHash(id, email, password string, salt []byte) (saltedPassword PassHash) {
	const (
		argonTimes   = 1
		argonMem     = 64 * 1024
		argonThreads = 4
		argonOut     = 32
	)
	saltedPassword.Salt = salt
	if len(saltedPassword.Salt) == 0 {
		saltedPassword.Salt = make([]byte, 128)
		_, err := rand.Read(saltedPassword.Salt)
		if err != nil {
			panic(err)
		}
	}
	hashedPassword := norm.NFD.Bytes([]byte(fmt.Sprintf("%s+%s+%s", id, email, password)))
	buf := make(
		[]byte,
		len([]byte(os.Getenv("PASSWORD_SECRET"))),
		len([]byte(os.Getenv("PASSWORD_SECRET")))+len(hashedPassword),
	)
	copy(buf, []byte(os.Getenv("PASSWORD_SECRET")))
	saltedPassword.Hash = argon2.IDKey(append(buf, hashedPassword...), saltedPassword.Salt, argonTimes, argonMem, argonThreads, argonOut)
	return saltedPassword
}

// EqualPassHash ...
func EqualPassHash(id, email, password string, saltedPassword PassHash) bool {
	comparedPassword := NewPassHash(id, email, password, saltedPassword.Salt)
	return hmac.Equal(saltedPassword.Hash, comparedPassword.Hash)
}
