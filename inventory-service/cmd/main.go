package main

import (
	"fmt"
	"inventory-service/infrastructure/db"
	"inventory-service/internal/handler"
	"inventory-service/internal/repository"
	"inventory-service/internal/usecase"
	"inventory-service/logger"
	pb "inventory-service/pb/inventory"
	"net"

	"google.golang.org/grpc"
)

func main() {
	logger.InitLogger()
	logger.Log.Info("🔄 Inventory service started")

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		logger.Log.Error(fmt.Sprintf("failed to listen: %v", err))
	}

	database := db.NewPostgres()
	productRepo := &repository.ProductRepository{DB: database}
	productUsecase := &usecase.ProductUsecase{Repo: productRepo}
	productHandler := &handler.ProductHandler{Usecase: productUsecase}

	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, productHandler)

	logger.Log.Info("Inventory gRPC service running on :50053")
	if err := grpcServer.Serve(listener); err != nil {
		logger.Log.Error(fmt.Sprintf("failed to serve: %v", err))
	}
}
