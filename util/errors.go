package util

import "errors"


// GetError ...
func GetError(err string) error {
	return errors.New(err)
}
