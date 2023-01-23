package main

import (
	"log"
	"net"

	"github.com/e-commerce-microservices/search-service/pb"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	es := esconn()

	grpcServer := grpc.NewServer()
	searchService := searchService{
		esconn: es,
	}
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
