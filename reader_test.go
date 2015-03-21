package akaibu

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	uncompressedHeader   = []byte{'A', 'K', 'A', 'I', 1, 0, 0, 0}
	uncompressedHeaderV2 = []byte{'A', 'K', 'A', 'I', 2, 0, 0, 0}
	uncompressed1B1      = []byte{'A', 'K', 'A', 'I', 1, 0, 0, 0, 1, 127}
	uncompressed1B2      = []byte{'A', 'K', 'A', 'I', 1, 0, 0, 0, 128, 1, 127}
	uncompressed1B3      = []byte{'A', 'K', 'A', 'I', 1, 0, 0, 0, 192, 0, 1, 127}
	uncompressed1B4      = []byte{'A', 'K', 'A', 'I', 1, 0, 0, 0, 224, 0, 0, 1, 127}
	uncompressed1B5      = []byte{'A', 'K', 'A', 'I', 1, 0, 0, 0, 240, 0, 0, 0, 1, 127}
)

func TestNewReader(t *testing.T) {
	// Initializing a new reader with less than header data fails with
	// io.ErrUnexpectedEOF.
	for l := 0; l < 7; l++ {
		_, err := NewReader(bytes.NewReader(uncompressedHeader[:l]))
		if err != io.ErrUnexpectedEOF {
			t.Fatalf("unexpected error creating reader with insufficient data: %v", err)
		}
	}

	// Initializing a new reader with unsupported version fails with
	// ErrUnsupportedVersion.
	_, err := NewReader(bytes.NewReader(uncompressedHeaderV2))
	if err != ErrUnsupportedVersion {
		t.Fatalf("unexpected error creating reader with unsupported version: %v", err)
	}

	// Initialize a new reader for an uncompressed, empty archive file
	// succeeds.
	r := bytes.NewReader(uncompressedHeader)
	ar, err := NewReader(r)
	if err != nil {
		t.Fatalf("unexpected error creating reader for uncompressed archive log file: %v", err)
	}

	if ar == nil {
		t.Fatalf("archive reader is unexpectedly nil")
	} else {
		ar.Close()
	}
}

func TestReaderRead(t *testing.T) {
	// Test reading from an uncompressed, empty archive file.
	ar, err := NewReader(bytes.NewReader(uncompressedHeader))
	if err != nil {
		t.Fatalf("unexpected error creating reader for uncompressed archive log file: %v", err)
	}
	defer ar.Close()

	d, err := ar.Read()
	if err != io.EOF {
		t.Fatalf("unexpected error reading record from empty archive log file: %v", err)
	}
	if d != nil {
		t.Fatalf("returned data from reading record is not nil: %v", d)
	}

	// Test reading with a size of 1 B encoded in 1-5 octets.
	for i, fixture := range [][]byte{
		uncompressed1B1,
		uncompressed1B2,
		uncompressed1B3,
		uncompressed1B4,
		uncompressed1B5,
	} {
		ar, err := NewReader(bytes.NewReader(fixture))
		if err != nil {
			t.Fatalf("unexpected error creating reader for uncompressed archive log file with 1 B record with length encoded in %d octets: %v", i+1, err)
		}
		defer ar.Close()

		d, err := ar.Read()
		if err != nil {
			t.Fatalf("unexpected error reading record from archive log file with 1 B record with length encoded in %d octets: %v", i+1, err)
		}

		if !bytes.Equal(d, []byte{127}) {
			t.Fatalf("unexpected data record read from archive log file with 1 B record with length encoded in %d octets: %v", i+1, d)
		}
	}
}

func TestReaderWithFixture(t *testing.T) {
	// Read out the list of available fixtures.
	//
	// Fixtures are assumed to be laid out like the samples generated in
	// akaibu-samples, ie., record 1 has a length of 100, record 2 a length of
	// 200 and so forth.
	fixturesFi, err := ioutil.ReadDir("fixtures")
	if err != nil {
		t.Logf("No fixtures available: %s", err)
		return
	}

	fixturePaths := make([]string, 0, len(fixturesFi))
	for _, fixtureFi := range fixturesFi {
		if !fixtureFi.Mode().IsRegular() {
			continue
		}

		fixturePaths = append(fixturePaths, filepath.Join("fixtures", fixtureFi.Name()))
	}

	// Run across the available fixtures.
	if len(fixturePaths) == 0 {
		t.Log("No fixtures available")
		return
	}

	for _, path := range fixturePaths {
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("unable to open test fixture %s for reading: %v", path, err)
		}
		defer f.Close()

		ar, err := NewReader(f)
		if err != nil {
			t.Fatalf("unexpected error creating reader for test fixture %s: %v", path, err)
		}
		defer ar.Close()

		i := 1
		for {
			r, err := ar.Read()
			t.Logf("Read() for %s#%d = %v, %v", path, i, r, err)

			if err == io.EOF {
				break
			} else if err != nil {
				t.Fatalf("unexpected error reading record from test fixture %s: %v", path, err)
			}

			if len(r) != i*100 {
				t.Fatalf("unexpected length of data record from test fixture %s: %d (expected %d)", path, len(r), i*100)
			}

			i++
		}
	}
}
