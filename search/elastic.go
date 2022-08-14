package search

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/ChrisCodeX/Event-Architecture-CQRS-Go/models"
	elastic "github.com/elastic/go-elasticsearch/v7"
)

type ElasticSearchRepository struct {
	client *elastic.Client
}

// Constructor
func NewElastic(url string) (*ElasticSearchRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}
	return &ElasticSearchRepository{client: client}, nil
}

// Close Method
func (e *ElasticSearchRepository) Close() {
	//
}

// Index feeds
func (e *ElasticSearchRepository) IndexFeed(ctx context.Context, feed models.Feed) error {
	body, _ := json.Marshal(feed)

	// ElasticSearch Index configuration
	_, err := e.client.Index(
		// Index name
		"feeds",
		// Reader that processes the body
		bytes.NewReader(body),
		// Id for documents
		e.client.Index.WithDocumentID(feed.Id),
		// Context for help to debug
		e.client.Index.WithContext(ctx),
		// Parameter that refresh the index
		e.client.Index.WithRefresh("wait_for"),
	)
	return err
}

// Search Feeds
func (e *ElasticSearchRepository) SearchFeed(ctx context.Context, query string) ([]models.Feed, error) {
	// Buffer
	var buf bytes.Buffer

	// Query that will be sent to ElasticSearch for the search
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			// Indicate to bring several documents
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title", "description"},
				// Smart search
				"fuzziness": 3,
				// Frequency of appearance of the query
				"cutoff_frequency": 0.0001,
			},
		},
	}
}
