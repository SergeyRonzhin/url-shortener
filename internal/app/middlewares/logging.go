package middlewares

import (
	"net/http"
	"time"
)

type (
	logRW struct {
		http.ResponseWriter
		status int
		size   int
	}
)

func (r *logRW) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size

	return size, err
}

func (r *logRW) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode
}

func (m *Middlewares) Logging(baseHandler http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		lw := logRW{
			ResponseWriter: w,
		}

		baseHandler.ServeHTTP(&lw, r)
		duration := time.Since(start)

		m.logger.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", lw.status,
			"size", lw.size,
		)
	}

	return http.HandlerFunc(logFn)
}
