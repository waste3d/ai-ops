package domain

import "time"

type Ticket struct {
	ID             string
	Source         string
	Payload        string
	Status         string
	AnalysisResult string
	CreatedAt      time.Time
}
