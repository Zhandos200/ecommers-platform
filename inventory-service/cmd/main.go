package main

import (
	"inventory-service/infrastructure/db"
	"inventory-service/internal/handler"
	"inventory-service/internal/repository"
	"inventory-service/internal/usecase"
	pb "inventory-service/pb/inventory"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	database := db.NewPostgres()
	productRepo := &repository.ProductRepository{DB: database}
	productUsecase := &usecase.ProductUsecase{Repo: productRepo}
	productHandler := &handler.ProductHandler{Usecase: productUsecase}

	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, productHandler)

	log.Println("Inventory gRPC service running on :50053")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
