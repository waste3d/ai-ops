package domain

import "time"

type Ticket struct {
	ID        string
	Source    string
	Payload   string
	CreatedAt time.Time
}
