package grpc

import (
	"context"

	authpb "github.com/waste3d/ai-ops/gen/go/auth"
	"github.com/waste3d/ai-ops/services/user_service/internal/application"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer
	useCase *application.UserUseCase
}

func NewServer(useCase *application.UserUseCase) *Server {
	return &Server{useCase: useCase}
}

func (s *Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.User, error) {
	user, err := s.useCase.Register(ctx, req.Username, req.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &authpb.User{
		Id:        user.ID,
		Username:  user.Username,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

func (s *Server) GetUserByUsername(ctx context.Context, req *authpb.GetUserByUsernameRequest) (*authpb.GetUserByUsernameResponse, error) {
	user, err := s.useCase.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &authpb.GetUserByUsernameResponse{
		Id:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}, nil
}
