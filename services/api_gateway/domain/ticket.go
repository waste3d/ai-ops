package domain

import "time"

type TicketView struct {
	ID             string    `json:"id"`
	Source         string    `json:"source"`
	Payload        string    `json:"payload"`
	Status         string    `json:"status"`
	AnalysisResult string    `json:"analysis_result"`
	CreatedAt      time.Time `json:"created_at"`
}
