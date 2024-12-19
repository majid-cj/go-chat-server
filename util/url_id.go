package util

import (
	"regexp"
)

var (
	profileIdRegex = regexp.MustCompile(`[0-7][0-9A-Z]{25}`)
)

// GetURLIds ...
func GetURLIds(path string) []string {
	return profileIdRegex.FindAllString(path, -1)
}
