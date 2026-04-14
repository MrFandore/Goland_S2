package httpapi

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func LoggingMiddleware(log *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)

		// Генерация request_id
		requestID := time.Now().UnixNano()
		reqIDStr := strconv.FormatInt(requestID, 10)

		// Установка заголовка ответа
		w.Header().Set("X-Request-Id", reqIDStr)

		// Добавление request_id в контекст запроса
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		r = r.WithContext(ctx)

		log.Info("incoming request",
			zap.Int64("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
		)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		log.Info("request completed",
			zap.Int64("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status_code", lrw.StatusCode()),
			zap.Duration("duration", duration),
		)
	})
}
