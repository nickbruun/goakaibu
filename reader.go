package akaibu

import (
	"bufio"
	"bytes"
	"io"
)

// Reader.
//
// Reads records from an Akaibu archival log file.
type Reader interface {
	io.Closer

	// Read a record.
	Read() ([]byte, error)
}

// Reader implementation.
type reader struct {
	rc io.ReadCloser
}

// New reader.
func NewReader(r io.Reader) (Reader, error) {
	// Wrap the reader in a buffered reader if it is not already.
	var br *bufio.Reader
	var ok bool

	if br, ok = r.(*bufio.Reader); !ok {
		br = bufio.NewReader(r)
	}

	// Read the header.
	header := make([]byte, 8)
	if _, err := io.ReadAtLeast(br, header, 8); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}

	if !bytes.Equal(header[:4], []byte{'A', 'K', 'A', 'I'}) {
		return nil, ErrInvalidArchive
	}

	if header[4] != 1 {
		return nil, ErrUnsupportedVersion
	}

	// Set up decompression if necessary.
	var comp Compression = Compression(header[5])
	if !comp.Valid() {
		return nil, ErrInvalidArchive
	}

	rc, err := comp.wrapReader(br)
	if err != nil {
		return nil, err
	}

	return &reader{
		rc: rc,
	}, nil
}

func (r *reader) Close() (err error) {
	if r.rc != nil {
		err = r.rc.Close()
		r.rc = nil
	}
	return
}

func (r *reader) Read() ([]byte, error) {
	// Read the prefix (at least one octet.)
	prefix := make([]byte, 1)
	_, err := io.ReadFull(r.rc, prefix)

	if err != nil {
		return nil, err
	}

	// Determine the size.
	var size uint64
	pi := uint8(prefix[0])

	if (pi & b10000000) == 0 {
		// One octet size indicator.
		size = uint64(pi & b01111111)
	} else if (pi & b11000000) == b10000000 {
		// Two octet size indicator.
		suffix := make([]byte, 1)
		if _, err = io.ReadFull(r.rc, suffix); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return nil, err
		}

		size = (uint64(pi&b00111111) << 8) | uint64(suffix[0])
	} else if (pi & b11100000) == b11000000 {
		// Three octet size indicator.
		suffix := make([]byte, 2)
		if _, err = io.ReadFull(r.rc, suffix); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return nil, err
		}

		size = (uint64(pi&b00011111) << 16) | (uint64(suffix[0]) << 8) | uint64(suffix[1])
	} else if (pi & b11110000) == b11100000 {
		// Four octet size indicator.
		suffix := make([]byte, 3)
		if _, err = io.ReadFull(r.rc, suffix); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return nil, err
		}

		size = (uint64(pi&b00001111) << 24) | (uint64(suffix[0]) << 16) | (uint64(suffix[1]) << 8) | uint64(suffix[2])
	} else if (pi & b11111000) == b11110000 {
		// Five octet size indicator.
		suffix := make([]byte, 4)
		if _, err = io.ReadFull(r.rc, suffix); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return nil, err
		}

		size = (uint64(pi&b00000111) << 32) | (uint64(suffix[0]) << 24) | (uint64(suffix[1]) << 16) | (uint64(suffix[2]) << 8) | uint64(suffix[3])
	} else {
		return nil, ErrInvalidArchive
	}

	// Read the residual data not contained in the prefix.
	data := make([]byte, size)
	if size == 0 {
		return data, nil
	}

	_, err = io.ReadFull(r.rc, data)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}

	return data, err
}
