package server

import (
	"context"
	"log"

	"github.com/MrFandore/Go_S2/Prac2/pkg/api/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
}

func New() *AuthServer {
	return &AuthServer{}
}

// Verify проверяет токен. Для демо считаем валидным токен "valid-token".
func (s *AuthServer) Verify(ctx context.Context, req *auth.VerifyRequest) (*auth.VerifyResponse, error) {
	log.Printf("Verify called with token: %s", req.Token)

	if req.Token == "" {
		return nil, status.Error(codes.Unauthenticated, "empty token")
	}

	// Простейшая валидация (замените на реальную логику)
	if req.Token != "valid-token" {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return &auth.VerifyResponse{
		Valid:   true,
		Subject: "user-123",
	}, nil
}
