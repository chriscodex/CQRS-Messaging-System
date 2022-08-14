package events

import (
	"context"

	"github.com/ChrisCodeX/Event-Architecture-CQRS-Go/models"
)

/* Abstract implementation of NATS */
type EventStore interface {
	Close()
	// Publish new feed event
	PublishCreatedFeed(ctx context.Context, feed *models.Feed) error
	// Subscribe to feed event
	SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error)
	// Callback when new feed is created
	OnCreateFeed(f func(CreatedFeedMessage)) error
}

var eventStore EventStore

/*Methods*/
func Close() {
	eventStore.Close()
}

func PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	return eventStore.PublishCreatedFeed(ctx, feed)
}

func SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	return eventStore.SubscribeCreatedFeed(ctx)
}

func OnCreateFeed(f func(CreatedFeedMessage)) error {
	return eventStore.OnCreateFeed(f)
}
