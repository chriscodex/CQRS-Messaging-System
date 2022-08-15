package main

import (
	"fmt"
	"log"

	"github.com/ChrisCodeX/Event-Architecture-CQRS-Go/database"
	"github.com/ChrisCodeX/Event-Architecture-CQRS-Go/repository"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

func main() {
	// Decode environment variables into Config struct
	var cfg Config
	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatalf("%v", err)
	}

	/*Database Connection*/
	// Get the url for create a postgres instance
	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

	// Create a new instance of postgres repository
	repo, err := database.NewPostgresRepository(addr)
	if err != nil {
		log.Fatalf("%v", err)
	}
	repository.SetRepository(repo)

}
