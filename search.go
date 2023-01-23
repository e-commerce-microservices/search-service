package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/e-commerce-microservices/search-service/pb"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/golang/protobuf/ptypes/empty"
)

type searchService struct {
	esconn *elasticsearch.Client
	pb.UnimplementedSearchServiceServer
}

const productIndex = "product"

func (srv *searchService) SearchProduct(ctx context.Context, req *pb.SearchProductRequest) (*pb.SearchProductResponse, error) {
	queryBody := getProductQueryBody(req.GetProductName())

	es := srv.esconn
	// Perform the search request
	res, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(productIndex),
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

	listProductName := make([]string, 0)
	listProductID := make([]int64, 0)
	limit := 5
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) < 5 {
		limit = len(hits)
	}
	// ID and document source for each hit.
	for i, hit := range hits {
		if i <= limit {
			source := hit.(map[string]interface{})["_source"].(map[string]interface{})
			listProductName = append(listProductName, source["product_name"].(string))
		}
		idStr := hit.(map[string]interface{})["_id"].(string)
		id, _ := strconv.ParseInt(idStr, 10, 64)
		listProductID = append(listProductID, id)
	}

	return &pb.SearchProductResponse{
		ListProductName: listProductName,
		ListProductId:   listProductID,
	}, nil
}

func getProductQueryBody(productName string) bytes.Buffer {
	var buf bytes.Buffer

	query := map[string]any{
		"query": map[string]any{
			"match": map[string]any{
				"product_name": productName,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("Error encoding query: %s\n", err)
	}

	return buf
}
func (srv *searchService) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*empty.Empty, error) {
	es := srv.esconn

	// Build the request body
	data, err := json.Marshal(struct {
		ProductName string `json:"product_name"`
	}{ProductName: req.GetProductName()})

	// Set up the request object
	indexRequest := esapi.IndexRequest{
		Index:      productIndex,
		DocumentID: fmt.Sprint(req.GetProductId()),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// Perform the request with the client
	res, err := indexRequest.Do(context.Background(), es)
	if err != nil {
		log.Printf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("[%s] Error indexing document ID=%d", res.Status(), req.GetProductId())
	}
	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("Error parsing the response body: %s", err)
	}

	return &empty.Empty{}, nil
}

func (srv *searchService) Ping(context.Context, *empty.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Message: "pong",
	}, nil
}
