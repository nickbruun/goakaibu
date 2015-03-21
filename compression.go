package akaibu

import (
	"code.google.com/p/snappy-go/snappy"
	"compress/zlib"
	"io"
	"io/ioutil"
)

// Compression.
type Compression uint8

const (
	// Uncompressed.
	UncompressedCompression Compression = iota

	// zlib compression.
	ZlibCompression

	// Snappy compression.
	SnappyCompression
)

// Test if the compression is valid.
func (c Compression) Valid() bool {
	return c <= SnappyCompression
}

// Wrap a reader for decompression.
func (c Compression) wrapReader(r io.Reader) (io.ReadCloser, error) {
	switch c {
	case UncompressedCompression:
		return ioutil.NopCloser(r), nil

	case ZlibCompression:
		return zlib.NewReader(r)

	case SnappyCompression:
		return ioutil.NopCloser(snappy.NewReader(r)), nil

	default:
		return nil, ErrInvalidCompression
	}
}
