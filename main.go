package main

import (
	"log"
	"net"

	"github.com/e-commerce-microservices/search-service/pb"
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
	seed(es)

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
