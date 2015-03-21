package akaibu

import (
	"encoding/binary"
)

const (
	b10000000 uint8 = 0x80
	b11000000 uint8 = 0xc0
	b11100000 uint8 = 0xe0
	b11110000 uint8 = 0xf0
	b11111000 uint8 = 0xf8

	b01111111 uint8 = 0x7f
	b00111111 uint8 = 0x3f
	b00011111 uint8 = 0x1f
	b00001111 uint8 = 0x0f
	b00000111 uint8 = 0x07
)

// Pack an integer value.
func PackInt(v uint64) ([]byte, error) {
	if v >= 34359738368 {
		return nil, ErrOutOfRange
	}

	if v < 128 {
		return []byte{byte(v)}, nil
	}

	d := make([]byte, 8)
	binary.BigEndian.PutUint64(d, v)

	if v < 16384 {
		return []byte{
			d[6] | b10000000,
			d[7],
		}, nil
	}

	if v < 2097152 {
		return []byte{
			d[5] | b11000000,
			d[6],
			d[7],
		}, nil
	}

	if v < 268435456 {
		return []byte{
			d[4] | b11100000,
			d[5],
			d[6],
			d[7],
		}, nil
	}

	return []byte{
		d[3] | b11110000,
		d[4],
		d[5],
		d[6],
		d[7],
	}, nil
}
