package akaibu

import (
	"io"
)

func writeFull(w io.Writer, p []byte) (err error) {
	var n int

	for len(p) > 0 {
		if n, err = w.Write(p); err != nil {
			break
		}

		p = p[n:]
	}

	return
}
