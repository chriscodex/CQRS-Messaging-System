package repository

import (
	"context"

	"github.com/ChrisCodeX/Event-Architecture-CQRS-Go/models"
)

type repository interface {
	Close()
	InsertFeed(ctx context.Context, feed *models.Feed) error
}
