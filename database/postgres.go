package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/ChrisCodeX/CQRS-Messaging-System/models"
)

type PostgresRepository struct {
	db *sql.DB
}

// Constructor
func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

// Close connection with database
func (repo *PostgresRepository) Close() {
	repo.db.Close()
}

// Insert feed
func (repo *PostgresRepository) InsertFeed(ctx context.Context, feed *models.Feed) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO feeds (id, title, description) VALUES ($1, $2, $3)", feed.Id, feed.Title, feed.Description)
	if err != nil {
		return err
	}
	return nil
}

// List of feeds
func (repo *PostgresRepository) ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	// Query
	rows, err := repo.db.QueryContext(ctx, "SELECT id, title, description, values FROM feeds")
	if err != nil {
		return nil, err
	}

	// Mapping
	var feeds []*models.Feed
	for rows.Next() {
		var feed *models.Feed
		if err := rows.Scan(&feed.Id, &feed.Title, &feed.Description, &feed.CreatedAt); err == nil {
			feeds = append(feeds, feed)
		}
	}
	return feeds, nil
}
