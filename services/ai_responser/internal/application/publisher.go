package application

import (
	"context"

	"github.com/waste3d/ai-ops/services/ai_responser/internal/domain"
)

type ResultPublisher interface {
	PublishAnalysisResult(ctx context.Context, result *domain.AnalysisResult) error
}
