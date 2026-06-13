package cryptex

import "errors"

var (
	ErrInvalidSize    = errors.New("invalid size")
	ErrInvalidPadding = errors.New("invalid padding") // unused
)
