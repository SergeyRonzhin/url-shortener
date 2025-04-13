package handlers

import (
	"net/http"
	"time"
)

type (
	rsData struct {
		status int
		size   int
	}

	logRW struct {
		http.ResponseWriter
		responseData *rsData
	}
)

func (r *logRW) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size

	return size, err
}

func (r *logRW) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func (h *HTTPHandler) Logging(baseHandler http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		lw := logRW{w, &rsData{}}
		baseHandler.ServeHTTP(&lw, r)

		duration := time.Since(start)

		h.logger.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", lw.responseData.status,
			"size", lw.responseData.size,
		)
	}

	return http.HandlerFunc(logFn)
}
