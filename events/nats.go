package events

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/ChrisCodeX/CQRS-Messaging-System/models"
	"github.com/nats-io/nats.go"
)

/*Concrete implementation of NATS*/

type NatsEventStore struct {
	// Connection
	conn *nats.Conn
	// Subscription which will use the feed struct to subscribe to an event
	feedCreatedSub *nats.Subscription
	//
	feedCreatedChan chan CreatedFeedMessage
}

/*Methods of NatsEventStore*/
// Constructor
func NewNats(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{conn: conn}, nil
}

// Close Method
func (n *NatsEventStore) Close() {
	// Verify if the connection exists
	if n.conn != nil {
		// Close the connection with the server
		n.conn.Close()
	}
	// Verify if the subscription exists
	if n.feedCreatedSub != nil {
		// Unsubscribe from event
		n.feedCreatedSub.Unsubscribe()
	}
	// Close channel of transmition of feeds
	close(n.feedCreatedChan)
}

// Function to encode messages to bytes
func (n *NatsEventStore) encodeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Method Publish Created Feed to services connected to NATS
func (n *NatsEventStore) PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	msg := CreatedFeedMessage{
		Id:          feed.Id,
		Title:       feed.Title,
		Description: feed.Description,
		CreatedAt:   feed.CreatedAt,
	}

	data, err := n.encodeMessage(msg)

	if err != nil {
		return err
	}

	// Publish the message with the corresponding data
	return n.conn.Publish(msg.Type(), data)
}

// Function to decode bytes to interface
func (n *NatsEventStore) decodeMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

// Closure to subscribe
func (n *NatsEventStore) OnCreateFeed(f func(CreatedFeedMessage)) (err error) {
	msg := CreatedFeedMessage{}

	n.feedCreatedSub, err = n.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		n.decodeMessage(m.Data, &msg)
		f(msg)
	})

	return
}

// Subscription Channel
func (n *NatsEventStore) SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	m := CreatedFeedMessage{}

	n.feedCreatedChan = make(chan CreatedFeedMessage, 64)

	ch := make(chan *nats.Msg, 64)

	var err error

	n.feedCreatedSub, err = n.conn.ChanSubscribe(m.Type(), ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				n.decodeMessage(msg.Data, &m)
				n.feedCreatedChan <- m
			}
		}
	}()
	return (<-chan CreatedFeedMessage)(n.feedCreatedChan), nil
}
