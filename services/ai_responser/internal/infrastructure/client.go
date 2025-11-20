package llm

import (
	"context"
	"strings"
	"time"

	"github.com/waste3d/ai-ops/services/ai_responser/internal/application"
)

type SimulatedClient struct{}

var _ application.Analyzer = (*SimulatedClient)(nil)

func NewSimulatedClient() *SimulatedClient {
	return &SimulatedClient{}
}

func (c *SimulatedClient) Analyze(ctx context.Context, payload string) (string, error) {
	time.Sleep(1 * time.Second)

	if strings.Contains(strings.ToLower(payload), "баз") {
		return "Проблема классифицирована как связанная с базой данных.", nil
	}
	if strings.Contains(strings.ToLower(payload), "диск") {
		return "Проблема классифицирована как связанная с дисковым пространством.", nil
	}
	return "Проблема классифицирована как связанная с другими компонентами системы.", nil
}
