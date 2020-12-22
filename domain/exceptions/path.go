package exceptions

import "errors"

var (
	ErrPathNotFound   = errors.New("path not found")
	ErrVertexNotFound = errors.New("invalid airport")
)
