package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/status"

	"github.com/MrFandore/Go_S2/Prac2/pkg/api/auth"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client auth.AuthServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &AuthClient{
		conn:   conn,
		client: auth.NewAuthServiceClient(conn),
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

// VerifyToken вызывает gRPC метод Verify с дедлайном.
func (c *AuthClient) VerifyToken(ctx context.Context, token string) (bool, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := c.client.Verify(ctx, &auth.VerifyRequest{Token: token})
	if err != nil {
		return false, "", err
	}
	return resp.Valid, resp.Subject, nil
}
