package persistence

import (
	"context"

	dbsql "database/sql"

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

func (r *PostgresTicketRepository) GetAll(ctx context.Context) ([]*domain.Ticket, error) {
	sql := `SELECT id, source, payload, status, analysis_result, created_at FROM tickets ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*domain.Ticket
	for rows.Next() {
		var t domain.Ticket
		var analysisResult dbsql.NullString

		if err := rows.Scan(&t.ID, &t.Source, &t.Payload, &t.Status, &analysisResult, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.AnalysisResult = analysisResult.String
		tickets = append(tickets, &t)
	}
	return tickets, nil
}

func (r *PostgresTicketRepository) GetByID(ctx context.Context, id string) (*domain.Ticket, error) {
	sql := `SELECT id, source, payload, status, analysis_result, created_at FROM tickets WHERE id = $1`
	row := r.db.QueryRow(ctx, sql, id)
	var t domain.Ticket
	var analysisResult dbsql.NullString

	if err := row.Scan(&t.ID, &t.Source, &t.Payload, &t.Status, &analysisResult, &t.CreatedAt); err != nil {
		return nil, err
	}
	t.AnalysisResult = analysisResult.String
	return &t, nil
}
