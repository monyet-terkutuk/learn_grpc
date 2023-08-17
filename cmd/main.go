package main

import (
	"log"
	"net"

	"learn_grpc/cmd/config"
	"learn_grpc/cmd/services"
	productPB "learn_grpc/pb/product"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	netListen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen %v", err.Error())
	}

	db := config.ConnectDB()

	grpcServer := grpc.NewServer()
	productService := services.ProductService{DB: db}
	productPB.RegisterProductServiceServer(grpcServer, &productService)

	log.Printf("Server running on %v", netListen.Addr())
	if err := grpcServer.Serve(netListen); err != nil {
		log.Fatalf("Failed to serve %v", err.Error())
	}
}
