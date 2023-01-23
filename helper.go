package main

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func esconn() *elasticsearch.Client {
	// Initialize a client with the config
	cfg := elasticsearch.Config{
		Addresses: []string{"http://elasticsearch-master-hl:9200"},
	}

	es, _ := elasticsearch.NewClient(cfg)

	// 1. Get cluster info
	//
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	return es
}
