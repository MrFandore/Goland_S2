package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter обёртка для захвата кода статуса.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware логирует каждый запрос: метод, путь, статус, длительность, request-id.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wr := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wr, r)
		duration := time.Since(start)
		requestID := GetRequestID(r.Context())
		log.Printf("[%s] %s %s %d %v", requestID, r.Method, r.URL.Path, wr.statusCode, duration)
	})
}
