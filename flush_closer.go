package akaibu

import (
	"io"
)

type flushWriter interface {
	io.Writer
	Flush() error
}

type flushCloser struct {
	flushWriter
}

func (c flushCloser) Close() error {
	return c.Flush()
}

func newFlushCloser(w flushWriter) io.WriteCloser {
	return flushCloser{w}
}
