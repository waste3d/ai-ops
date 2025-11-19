package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/waste3d/ai-ops/services/auditor/internal/domain"
)

type PostgresTicketRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTicketRepository(db *pgxpool.Pool) *PostgresTicketRepository {
	return &PostgresTicketRepository{db: db}
}

func (r *PostgresTicketRepository) Save(ctx context.Context, ticket *domain.Ticket) error {
	sql := `INSERT INTO tickets (id, source, payload, status, created_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING`
	_, err := r.db.Exec(ctx, sql, ticket.ID, ticket.Source, ticket.Payload, ticket.Status, ticket.CreatedAt)
	return err
}

func (r *PostgresTicketRepository) Update(ctx context.Context, ticketID string, status string, result string) error {
	sql := `UPDATE tickets SET status = $1, analysis_result = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, sql, status, result, ticketID)
	return err
}
