package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/waste3d/ai-ops/services/user_service/internal/domain"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}
func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	sql := `INSERT INTO users (username, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, sql, user.Username, user.PasswordHash, user.CreatedAt).Scan(&user.ID, &user.CreatedAt)
	return user, err
}

func (r *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	sql := `SELECT id, username, password_hash, created_at FROM users WHERE username = $1`
	row := r.db.QueryRow(ctx, sql, username)
	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	return &u, err
}
