package llm

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/waste3d/ai-ops/services/ai_responser/internal/application"
	"google.golang.org/api/option"
)

type GoogleAIClient struct {
	client *genai.GenerativeModel
}

var _ application.Analyzer = (*GoogleAIClient)(nil)

func NewGoogleAIClient(apiKey string) (*GoogleAIClient, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("could not create Google AI client: %w", err)
	}

	model := client.GenerativeModel("gemini-2.5-flash")
	return &GoogleAIClient{client: model}, nil
}

func (c *GoogleAIClient) Analyze(ctx context.Context, payload string) (string, error) {
	prompt := fmt.Sprintf(`
	You are an experienced IT Operations engineer.
	Analyze the following problem description from a monitoring system and suggest a likely root cause or a next step for diagnostics.
	Be concise and clear.
	Answer on russian language only. Tell how to fix the problem or what to do next.

	Problem: "%s"`, payload)

	resp, err := c.client.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("Google AI content generation error: %w", err)
	}

	var result string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					result += string(txt)
				}
			}
		}
	}

	if result == "" {
		return "", fmt.Errorf("received empty response from Google AI")
	}

	return result, nil
}
