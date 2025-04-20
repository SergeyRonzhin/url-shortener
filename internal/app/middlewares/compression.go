package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (m *Middlewares) Compression(next http.Handler) http.Handler {
	compressFn := func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzReader, err := gzip.NewReader(r.Body)

		if err != nil {
			http.Error(w, "invalid gzip request", http.StatusBadRequest)
			return
		}

		defer gzReader.Close()
		r.Body = gzReader

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzWriter := gzip.NewWriter(w)
		defer gzWriter.Close()

		w.Header().Set("Content-Encoding", "gzip")
		r.Header.Set("Content-Type", "text/plain")

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzWriter}, r)
	}

	return http.HandlerFunc(compressFn)
}
