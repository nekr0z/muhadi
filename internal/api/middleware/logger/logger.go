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
			}

			next.ServeHTTP(&lw, r)

			duration := time.Since(start)

			log.Info("response served",
				zap.String("uri", uri),
				zap.String("method", method),
				zap.Duration("duration", duration),
				zap.Int("status", lw.status),
			)
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.WriteHeader(http.StatusOK)
	}

	size, err := r.ResponseWriter.Write(b)
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode
}
