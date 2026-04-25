package sse

import (
	"bytes"
	"fmt"
	"net/http"
)

type Writer struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

func NewWriter(w http.ResponseWriter) (*Writer, bool) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, false
	}
	return &Writer{w: w, flusher: flusher}, true
}

func (s *Writer) Write(p []byte) (n int, err error) {
	remaining := p
	for len(remaining) > 0 {
		var line []byte
		if i := bytes.IndexByte(remaining, '\n'); i >= 0 {
			line = remaining[:i]
			remaining = remaining[i+1:]
		} else {
			line = remaining
			remaining = nil
		}
		if len(line) == 0 {
			continue
		}
		if _, fmtError := fmt.Fprintf(s.w, "data: %s\n\n", line); fmtError != nil {
			return 0, fmtError
		}
	}
	s.flusher.Flush()
	return len(p), nil
}

func (s *Writer) WriteError(err error) {
	_, fmtError := fmt.Fprintf(s.w, "event: error\ndata: %s\n\n", err.Error())
	if fmtError != nil {
		return
	}
	s.flusher.Flush()
}
