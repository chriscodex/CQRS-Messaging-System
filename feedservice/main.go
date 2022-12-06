package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ChrisCodeX/CQRS-Messaging-System/database"
	"github.com/ChrisCodeX/CQRS-Messaging-System/events"
	"github.com/ChrisCodeX/CQRS-Messaging-System/repository"
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

	/*NATS Connection*/
	// Instance of NATS
	n, err := events.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		log.Fatal(err)
	}
	events.SetEventStore(n)

	// Close NATS Connection
	defer events.Close()

	// Listen and Serve
	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
