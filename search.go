package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/e-commerce-microservices/search-service/pb"
	"github.com/elastic/go-elasticsearch/v8"
)

type searchService struct {
	esconn *elasticsearch.Client
	pb.UnimplementedSearchServiceServer
}

func (srv *searchService) SearchProduct(ctx context.Context, req *pb.SearchProductRequest) (*pb.SearchProductResponse, error) {
	queryBody := getProductQueryBody(req.GetProductName())

	es := srv.esconn
	// Perform the search request
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
		es.Search.WithBody(&queryBody),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("Error getting response: %s\n", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("Error parsing the response body: %s", err)
		}
		return nil, fmt.Errorf("Error: %s", e["error"])
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("Error parsing the response body: %s", err)
	}

	listProduct := make([]*pb.SearchProductResponse_ProductInfo, 0)
	// ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s\n", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	return &pb.SearchProductResponse{
		ListProduct: listProduct,
	}, nil
}

func getProductQueryBody(productName string) bytes.Buffer {
	var buf bytes.Buffer

	query := map[string]any{
		"query": map[string]any{
			"match": map[string]any{
				"title": "test",
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("Error encoding query: %s\n", err)
	}

	return buf
}
