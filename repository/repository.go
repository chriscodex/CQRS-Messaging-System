package repository

import (
	"context"

	"github.com/ChrisCodeX/Event-Architecture-CQRS-Go/models"
)

type Repository interface {
	Close()
	InsertFeed(ctx context.Context, feed *models.Feed) error
	ListFeeds(ctx context.Context) ([]*models.Feed, error)
}

var repository Repository

func SetRepository(r Repository) {
	repository = r
}
