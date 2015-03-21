package akaibu

import (
	"bytes"
	"testing"
)

func TestPackInt(t *testing.T) {
	for _, fixture := range []struct {
		Value    uint64
		Expected []byte
	}{
		{0, []byte{0}},
		{127, []byte{127}},
		{128, []byte{128, 128}},
		{16383, []byte{191, 255}},
		{16384, []byte{192, 64, 0}},
		{2097151, []byte{223, 255, 255}},
		{2097152, []byte{224, 32, 0, 0}},
		{268435455, []byte{239, 255, 255, 255}},
		{268435456, []byte{240, 16, 0, 0, 0}},
		{34359738367, []byte{247, 255, 255, 255, 255}},
		{34359738368, nil},
	} {
		actual, err := PackInt(fixture.Value)

		if fixture.Value >= 34359738368 {
			if err != ErrOutOfRange {
				t.Fatalf("unexpected error for packing integer out of range: %v", err)
			}
			if actual != nil {
				t.Fatalf("unexpected value for packing integer out of range: %v", actual)
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error packing %d: %v", fixture.Value, err)
			}

			if !bytes.Equal(actual, fixture.Expected) {
				t.Fatalf("unexpected packed value for %d: %v, expected %v", fixture.Value, actual, fixture.Expected)
			}
		}
	}
}
