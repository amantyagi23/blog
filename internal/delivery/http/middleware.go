package http

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"usermanagement/internal/infrastructure/logger"
)

// LoggingMiddleware logs HTTP requests.
func LoggingMiddleware(logger *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			ww := &responseWriter{w, http.StatusOK}

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			logger.Info("http request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.statusCode),
				zap.Duration("duration", duration),
				zap.String("request_id", middleware.GetReqID(r.Context())),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}