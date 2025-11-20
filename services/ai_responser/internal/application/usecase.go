package application

import (
	"context"

	"github.com/waste3d/ai-ops/services/ai_responser/internal/domain"
)

type Analyzer interface {
	Analyze(ctx context.Context, payload string) (string, error)
}

type AnalysisUseCase struct {
	analyzer  Analyzer
	publisher ResultPublisher
}

func NewAnalysisUseCase(analyzer Analyzer, publisher ResultPublisher) *AnalysisUseCase {
	return &AnalysisUseCase{analyzer: analyzer, publisher: publisher}
}

func (uc *AnalysisUseCase) AnalyzeTicket(ctx context.Context, ticket *domain.Ticket) error {
	analysisText, err := uc.analyzer.Analyze(ctx, ticket.Payload)
	if err != nil {
		return err
	}

	result := &domain.AnalysisResult{
		TicketID: ticket.ID,
		Result:   analysisText,
	}

	return uc.publisher.PublishAnalysisResult(ctx, result)
}
