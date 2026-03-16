package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MrFandore/Go_S2/Prac2/services/tasks/internal/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	UserSubjectKey contextKey = "subject"
)

func AuthMiddleware(authClient *client.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			token := parts[1]

			valid, subject, err := authClient.VerifyToken(r.Context(), token)
			if err != nil {
				if st, ok := status.FromError(err); ok {
					switch st.Code() {
					case codes.Unauthenticated:
						http.Error(w, "invalid token", http.StatusUnauthorized)
					case codes.DeadlineExceeded:
						http.Error(w, "auth service timeout", http.StatusGatewayTimeout)
					default:
						http.Error(w, "auth service error", http.StatusServiceUnavailable)
					}
				} else {
					http.Error(w, "auth service unavailable", http.StatusServiceUnavailable)
				}
				return
			}
			if !valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserSubjectKey, subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
