package events

import "time"

// Dependency inversion interface
type Message interface {
	Type() string
}

// Struct that will be transmitted through the NATS
type CreatedFeedMessage struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
