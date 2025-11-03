package prometheus

import (
	"net/http"
	"time"
)

// writer wraps http.ResponseWriter to track response metrics
type writer struct {
	w http.ResponseWriter

	// Response tracking
	code         int
	bytesWritten int64

	// Timing tracking
	arrivalTime  time.Time
	queueStart   time.Time
	processStart time.Time

	// Request tracking
	requestSize int64
}

func (w *writer) Flush() {
	if fl, ok := w.w.(http.Flusher); ok {
		fl.Flush()
	}
}

func (w *writer) WriteHeader(code int) {
	if w.code == -1 {
		w.code = code
	}
	w.w.WriteHeader(code)
}

func (w *writer) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	w.bytesWritten += int64(n)
	return n, err
}

func (w *writer) Header() http.Header {
	return w.w.Header()
}

// reset prepares the writer for reuse from the pool
func (w *writer) reset() {
	w.code = -1
	w.bytesWritten = 0
	w.arrivalTime = time.Time{}
	w.queueStart = time.Time{}
	w.processStart = time.Time{}
	w.requestSize = 0
	w.w = nil
}
