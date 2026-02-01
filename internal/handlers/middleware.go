package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rodnney/transaction-processor/pkg/logger"
)

func LoggingMiddleware(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}
			w.Header().Set("X-Request-ID", requestID)

			next.ServeHTTP(w, r)

			l.Info("Request processed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start).String(),
				"request_id", requestID,
			)
		})
	}
}

func RecoveryMiddleware(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					l.Error("Panic recovered", "error", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
