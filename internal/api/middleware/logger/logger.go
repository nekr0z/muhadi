package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Log(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			uri := r.RequestURI
			method := r.Method

			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   &responseData{},
			}

			next.ServeHTTP(&lw, r)

			duration := time.Since(start)

			log.Info("request",
				zap.String("uri", uri),
				zap.String("method", method),
				zap.Duration("duration", duration),
				zap.Int("status", lw.responseData.status),
				zap.Int("size", lw.responseData.size),
			)
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	if r.responseData.status == 0 {
		r.WriteHeader(http.StatusOK)
	}

	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

type responseData struct {
	status int
	size   int
}
