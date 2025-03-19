package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Logging(log *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			lrw := &loggingResponseWriter{w, http.StatusOK}

			next.ServeHTTP(lrw, r)

			duration := time.Since(startTime)

			log.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"HTTP request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.Int("status_code", lrw.statusCode),
				slog.Duration("duration", duration),
			)
		})
	}
}
