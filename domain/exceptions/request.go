package exceptions

import "errors"

var (
	ErrBadParameters = errors.New("invalid parameters sent in request")
)
