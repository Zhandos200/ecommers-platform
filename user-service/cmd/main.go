package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"user-service/infrastructure"
	"user-service/internal/handler"
	"user-service/logger"
	pb "user-service/pb/user"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	logger.InitLogger()
	logger.Log.Info("🔄 User service started")
	// Connect to DB
	db := infrastructure.NewPostgres()

	// Create a gRPC server
	grpcServer := grpc.NewServer()

	// Create handler that implements pb.UserServiceServer
	userHandler := handler.NewUserHandler(db)

	// Register the handler with gRPC
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	// Start listening
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to listen: %v", err))
	}

	log.Println("User gRPC server running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to serve gRPC server: %v", err))
	}
}
