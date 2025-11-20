package grpc_client

import (
	"context"

	authpb "github.com/waste3d/ai-ops/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client authpb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(ctx context.Context, addr string) (*UserClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := authpb.NewAuthServiceClient(conn)
	return &UserClient{client: client, conn: conn}, nil
}

func (c *UserClient) Close() {
	c.conn.Close()
}

func (c *UserClient) Register(ctx context.Context, username, passwordHash string) (*authpb.User, error) {
	req := &authpb.RegisterRequest{
		Username:     username,
		PasswordHash: passwordHash,
	}
	return c.client.Register(ctx, req)
}

func (c *UserClient) GetUserByUsername(ctx context.Context, username string) (*authpb.GetUserByUsernameResponse, error) {
	req := &authpb.GetUserByUsernameRequest{
		Username: username,
	}
	return c.client.GetUserByUsername(ctx, req)
}
