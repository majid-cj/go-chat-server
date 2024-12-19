package util

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// ULID ...
func ULID() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	ULID, _ := ulid.New(ms, entropy)
	return ULID.String()
}
