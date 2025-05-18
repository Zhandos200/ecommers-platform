package main

import (
	"fmt"
	"net"
	"order-service/infrastructure/db"
	"order-service/internal/handler"
	"order-service/internal/repository"
	"order-service/internal/usecase"
	"order-service/logger"
	pbInventory "order-service/pb/inventory"
	pb "order-service/pb/order"

	"net/http"
	"order-service/internal/nats"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	logger.InitLogger()
	logger.Log.Info("🔄 Order service started")

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to listen on port 50052: %v", err))
	}

	conn, err := grpc.Dial("inventory-service:50053", grpc.WithInsecure())
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Could not connect to Inventory Service: %v", err))
	}

	natsPublisher := nats.NewNatsPublisher("nats://nats:4222")

	inventoryClient := pbInventory.NewInventoryServiceClient(conn)

	database := db.NewPostgres()
	orderRepo := &repository.OrderRepository{DB: database}
	orderUsecase := &usecase.OrderUsecase{Repo: orderRepo}
	orderHandler := handler.NewOrderHandler(orderUsecase, inventoryClient, natsPublisher)

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderHandler)

	logger.Log.Info("Order gRPC service running on :50052")
	if err := grpcServer.Serve(listener); err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to serve gRPC server: %v", err))
	}
}
