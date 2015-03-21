package akaibu

import (
	"bufio"
	"code.google.com/p/snappy-go/snappy"
	"compress/zlib"
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

func (w *writer) Close() (err error) {
	if w.w != nil {
		err = w.w.Close()
		w.w = nil
	}
	return
}

func (w *writer) Write(p []byte) (err error) {
	return writeRecord(w.w, p)
}

// Snappy-specific writer implementation.
type snappyWriter struct {
	bw  *bufio.Writer
	sw  *snappy.Writer
	sbw *bufio.Writer
}

func (w *snappyWriter) Close() (err error) {
	if w.sbw != nil {
		err = w.sbw.Flush()
		err = w.bw.Flush()
		w.sbw = nil
		w.bw = nil
		w.sw = nil
	}
	return
}

func (w *snappyWriter) Write(p []byte) (err error) {
	return writeRecord(w.sbw, p)
}

// Write a header.
func writeHeader(w io.Writer, c Compression) error {
	return writeFull(w, []byte{'A', 'K', 'A', 'I', 1, byte(c), 0, 0})
}

// Write a record.
func writeRecord(w io.Writer, p []byte) (err error) {
	var sizeData []byte
	if sizeData, err = PackInt(uint64(len(p))); err != nil {
		return
	}

	if err = writeFull(w, sizeData); err != nil {
		return
	}

	if len(p) == 0 {
		return nil
	}

	return writeFull(w, p)
}

// New uncompressed writer.
func NewUncompressedWriter(w io.Writer) (Writer, error) {
	var bw *bufio.Writer
	var ok bool
	if bw, ok = w.(*bufio.Writer); !ok {
		bw = bufio.NewWriter(w)
	}

	// Attempt to write the header.
	if err := writeHeader(bw, UncompressedCompression); err != nil {
		return nil, err
	}

	return &writer{newFlushCloser(bw)}, nil
}

// New zlib-compressed writer.
//
// The level has the same function as in compress/zlib.
func NewZlibCompressedWriter(w io.Writer, level int) (Writer, error) {
	// Attempt to write the header.
	if err := writeHeader(w, ZlibCompression); err != nil {
		return nil, err
	}

	// Set up the compression writer.
	wc, err := zlib.NewWriterLevel(w, level)
	if err != nil {
		return nil, err
	}

	return &writer{wc}, nil
}

// New Snappy-compressed writer.
func NewSnappyCompressedWriter(w io.Writer) (Writer, error) {
	var bw *bufio.Writer
	var ok bool
	if bw, ok = w.(*bufio.Writer); !ok {
		bw = bufio.NewWriter(w)
	}

	// Attempt to write the header.
	if err := writeHeader(w, SnappyCompression); err != nil {
		return nil, err
	}

	// Set up the compression writer and write an empty slice to force the
	// header to be written.
	sw := snappy.NewWriter(bw)
	if _, err := sw.Write([]byte{}); err != nil {
		return nil, err
	}

	return &snappyWriter{bw, sw, bufio.NewWriterSize(sw, 1<<16)}, nil
}
