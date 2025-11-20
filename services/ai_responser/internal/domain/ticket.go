package domain

type Ticket struct {
	ID      string
	Payload string
}

type AnalysisResult struct {
	TicketID string
	Result   string
}
