package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "requestid"

// RequestIDMiddleware добавляет X-Request-ID в контекст и ответ,
// если заголовок отсутствует – генерирует новый UUID.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("x-request-id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		w.Header().Set("x-request-id", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID извлекает request-id из контекста.
func GetRequestID(ctx context.Context) string {
	val := ctx.Value(RequestIDKey)
	if val == nil {
		return ""
	}
	return val.(string)
}
