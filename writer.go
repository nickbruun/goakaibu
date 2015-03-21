package akaibu

import (
	"io"
)

// Writer.
type Writer interface {
	io.Closer

	// Write a record.
	Write(p []byte) error
}
