package database

import (
	"database/sql"

	_ "github.com/lib/pq"
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
