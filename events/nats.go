package events

import "github.com/nats-io/nats.go"

/*Concrete implementation of NATS*/

type NatsEventStore struct {
	// Connection
	conn *nats.Conn
	// Subscription which will use the feed struct to subscribe to an event
	feedCreatedSub *nats.Subscription
	//
	feedCreatedChan chan CreatedFeedMessage
}

// Constructor
func NewNats(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{conn: conn}, nil
}
