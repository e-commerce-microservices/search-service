package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/e-commerce-microservices/search-service/pb"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	// postgres driver
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Println("connect es")
	es := esconn()
	fmt.Println("es: ", es)

	log.Println("create grpc server")
	grpcServer := grpc.NewServer()
	searchService := searchService{
		esconn: es,
	}
	log.Println("register grpc server")
	pb.RegisterSearchServiceServer(grpcServer, &searchService)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}
}

func esconn() *elasticsearch.Client {
	// Initialize a client with the config
	cfg := elasticsearch.Config{
		Addresses: []string{"http://host.minikube.internal:9200"},
		// Addresses: []string{"http://localhost:9200"},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println("es error: ", err)
	}

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

func seed(es *elasticsearch.Client) {

	pgDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"192.168.49.2", "30012", "admin", "admin", "product",
	)
	productDB, err := sql.Open("postgres", pgDSN)
	defer productDB.Close()

	rows, err := productDB.Query("SELECT id,name FROM product")
	if err != nil {
		log.Println(err)
	}
	var productName string
	var productID int64
	for rows.Next() {
		err := rows.Scan(&productID, &productName)
		if err != nil {
			log.Println(err)
			continue
		}

		data, err := json.Marshal(struct {
			ProductName string `json:"product_name"`
		}{ProductName: productName})

		// Set up the request object
		indexRequest := esapi.IndexRequest{
			Index:      productIndex,
			DocumentID: fmt.Sprint(productID),
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
			log.Printf("[%s] Error indexing document ID=%d\n", res.Status(), productID)
		}
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s\n", err)
		} else {
			log.Println("Insert document ID success")
		}
	}

}
