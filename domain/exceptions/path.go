package exceptions

import "errors"

var (
	ErrPathNotFound       = errors.New("path not found")
	ErrVertexNotFound     = errors.New("invalid airport")
	ErrRouteAlreadyExists = errors.New("route already exists")

	ErrInvalidInputFile = errors.New("invalid input file")
)
