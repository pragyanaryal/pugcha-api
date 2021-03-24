package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type ResponseWriter struct {
	gw *gzip.Writer
	rw http.ResponseWriter
}

func (wr *ResponseWriter) Header() http.Header {
	return wr.rw.Header()
}

func (wr *ResponseWriter) Write(bytes []byte) (int, error) {
	return wr.gw.Write(bytes)
}

func (wr *ResponseWriter) WriteHeader(statusCode int) {
	wr.rw.WriteHeader(statusCode)
}

func (wr *ResponseWriter) Flush() {
	_ = wr.gw.Flush()
	_ = wr.gw.Close()
}

func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	gw := gzip.NewWriter(rw)
	return &ResponseWriter{rw: rw, gw: gw}
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			rw.Header().Set("Content-Encoding", "gzip")
			wrw := NewResponseWriter(rw)

			next.ServeHTTP(wrw, r)
			defer wrw.gw.Close()

			return
		}
		next.ServeHTTP(rw, r)
	})
}