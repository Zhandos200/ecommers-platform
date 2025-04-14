package main

import (
	"log"
	"net"
	"order-service/infrastructure/db"
	"order-service/internal/handler"
	"order-service/internal/repository"
	"order-service/internal/usecase"
	pbInventory "order-service/pb/inventory"
	pb "order-service/pb/order"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	conn, err := grpc.Dial("inventory-service:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to Inventory Service: %v", err)
	}

	inventoryClient := pbInventory.NewInventoryServiceClient(conn)

	database := db.NewPostgres()
	orderRepo := &repository.OrderRepository{DB: database}
	orderUsecase := &usecase.OrderUsecase{Repo: orderRepo}
	orderHandler := handler.NewOrderHandler(orderUsecase, inventoryClient)

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderHandler)

	log.Println("✅ Order gRPC service running on :50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
