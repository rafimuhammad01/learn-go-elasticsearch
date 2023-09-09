package main

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticSearch(cloudID, apiKey string) (*elasticsearch.TypedClient, error) {
	cfg := elasticsearch.Config{
		CloudID: cloudID,
		APIKey:  apiKey,
	}

	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		log.Fatalf("Error creating elasticsearch typed client: %s", err)
	}

	ctx := context.Background()
	isSuccess, err := es.Ping().Do(ctx)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	if isSuccess {
		log.Printf("ping to elasticsearch success")
	}

	return es, nil
}
