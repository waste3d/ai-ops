package application

import (
	"context"

	"github.com/waste3d/ai-ops/services/user_service/internal/domain"
)

type UserUseCase struct {
	repo UserRepository
}

func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) Register(ctx context.Context, username, passwordHash string) (*domain.User, error) {
	user := &domain.User{Username: username, PasswordHash: passwordHash}
	return uc.repo.Create(ctx, user)
}

func (uc *UserUseCase) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return uc.repo.GetByUsername(ctx, username)
}
