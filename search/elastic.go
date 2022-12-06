package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/ChrisCodeX/CQRS-Messaging-System/models"
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

	// Encode searchQuery into buff
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	/* Search in ElasticSearch*/
	// Search Configuration
	res, err := e.client.Search(
		// Context for help to debug
		e.client.Search.WithContext(ctx),
		// Indicate index
		e.client.Search.WithIndex("feeds"),
		// Send the buf(query encoded)
		e.client.Search.WithBody(&buf),
		// Get the total hits
		e.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}

	var feeds []models.Feed
	// Close the response body
	defer func() {
		if err := res.Body.Close(); err != nil {
			feeds = nil
		}
	}()

	// Check if error exists
	if res.IsError() {
		return nil, errors.New(res.String())
	}

	// Decode the response body of the search into eRes
	var eRes map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&eRes); err != nil {
		return nil, err
	}

	// Get the hits from the search
	for _, hit := range eRes["hits"].(map[string]interface{})["hits"].([]interface{}) {
		feed := models.Feed{}
		source := hit.(map[string]interface{})["_source"]

		// marshal source into bytes
		marshal, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}

		// unamarshal the marshan into feed
		if err := json.Unmarshal(marshal, &feed); err == nil {
			feeds = append(feeds, feed)
		}
	}
	return feeds, nil
}
