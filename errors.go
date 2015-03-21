package akaibu

import (
	"errors"
)

var (
	// Invalid archive error.
	ErrInvalidArchive = errors.New("invalid archive")

	// Unsupported version error.
	ErrUnsupportedVersion = errors.New("unsupported archive version")

	// Invalid compression error.
	ErrInvalidCompression = errors.New("invalid compression")

	// Out of range error.
	ErrOutOfRange = errors.New("out of range")
)
