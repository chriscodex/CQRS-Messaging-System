package main

import (
	"encoding/json"
	"net/http"
	"time"

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

}
