package akaibu

import (
	"bufio"
	"io"
)

// Writer.
type Writer interface {
	io.Closer

	// Write a record.
	Write(p []byte) error
}

// Writer implementation.
type writer struct {
	w io.WriteCloser
}

// New writer from io.WriteCloser.
func newWriter(w io.WriteCloser, c Compression) (Writer, error) {
	wr := &writer{
		w: w,
	}

	header := []byte{'A', 'K', 'A', 'I', 1, byte(c), 0, 0}
	if err := wr.writeAll(header); err != nil {
		wr.Close()
		return nil, err
	}

	return wr, nil
}

// New uncompressed writer.
func NewUncompressedWriter(w io.Writer) (Writer, error) {
	var bw *bufio.Writer
	var ok bool
	if bw, ok = w.(*bufio.Writer); !ok {
		bw = bufio.NewWriter(w)
	}

	return newWriter(newFlushCloser(bw), UncompressedCompression)
}

func (w *writer) Close() (err error) {
	if w.w != nil {
		err = w.w.Close()
		w.w = nil
	}
	return
}

func (w *writer) Write(p []byte) (err error) {
	var sizeData []byte
	if sizeData, err = PackInt(uint64(len(p))); err != nil {
		return
	}

	if err = w.writeAll(sizeData); err != nil {
		return
	}

	if len(p) == 0 {
		return nil
	}

	return w.writeAll(p)
}

func (w *writer) writeAll(p []byte) (err error) {
	var n int

	for len(p) > 0 {
		if n, err = w.w.Write(p); err != nil {
			break
		}

		p = p[n:]
	}

	return
}
