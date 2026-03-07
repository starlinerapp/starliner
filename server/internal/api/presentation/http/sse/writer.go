package sse

import (
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
	fmt.Fprintf(s.w, "data: %s\n\n", p)
	s.flusher.Flush()
	return len(p), nil
}

func (s *Writer) WriteError(err error) {
	fmt.Fprintf(s.w, "event: error\ndata: %s\n\n", err.Error())
	s.flusher.Flush()
}
