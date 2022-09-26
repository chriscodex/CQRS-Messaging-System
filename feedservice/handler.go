package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ChrisCodeX/CQRS-Messaging-System/events"
	"github.com/ChrisCodeX/CQRS-Messaging-System/models"
	"github.com/ChrisCodeX/CQRS-Messaging-System/repository"
	"github.com/segmentio/ksuid"
)

type createFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createFeedHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the request from server
	var req createFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Time of feed created
	createdAt := time.Now().UTC()

	// Generate id
	id, err := ksuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Feed wich be inserted in database
	feed := models.Feed{
		Id:          id.String(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   createdAt,
	}

	// Insert feed into database
	if err := repository.InsertFeed(r.Context(), &feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Publish event created feed
	if err := events.PublishCreatedFeed(r.Context(), &feed); err != nil {
		log.Printf("failed to publish created feed event: %v", err)
	}

	// Send feed to writer (client)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}
