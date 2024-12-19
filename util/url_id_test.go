package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetURLIds(t *testing.T) {
	ULID1 := ULID()
	ULID2 := ULID()

	profileIds := GetURLIds(fmt.Sprintf("/api/v1/profile/%s/%s", ULID1, ULID2))
	assert.Equal(t, profileIds[0], ULID1)
	assert.Equal(t, profileIds[1], ULID2)
}
