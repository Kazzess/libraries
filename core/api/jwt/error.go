package jwt

import (
	"git.adapticode.com/libraries/golang/errors"
)

var ErrBadToken = errors.New("malformed token")
