package trafficshaping

import (
	"io"

	"github.com/dty1er/traffic-shaping-io/bucket"
)

type Writer struct {
	w      io.Writer
	bucket bucket.Bucket
}

func (w *Writer) Write(p []byte) (int, error) {
	return w.w.Write(p)
}
