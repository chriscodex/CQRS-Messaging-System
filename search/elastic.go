package search

import (
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
