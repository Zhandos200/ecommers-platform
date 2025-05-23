package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"user-service/infrastructure"

	"user-service/internal/handler"
	"user-service/internal/repository"
	"user-service/internal/usecase"
	"user-service/logger"
	userpb "user-service/pb/user"
)

func main() {
	// 1) Prometheus metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Fatalf("metrics server error: %v", err)
		}
	}()

	// 2) Init logger
	logger.InitLogger()
	logger.Log.Info("🔄 User service started")

	// 3) Connect to Postgres
	db := infrastructure.NewPostgres()

	// 4) Конструируем слои приложения:
	// 4.1) Репозиторий
	repo := repository.NewUserRepository(db)
	// 4.2) Mailer для отправки писем
	mailer := usecase.NewSMTPMailer()
	// 4.3) Usecase (repo + mailer)
	uc := usecase.NewUserUsecase(repo, mailer)
	// 4.4) Handler (gRPC)
	userHandler := handler.NewUserHandler(uc)

	// 5) Запускаем gRPC-сервер
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to listen: %v", err))
		return
	}
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, userHandler)

	logger.Log.Info("🚀 User gRPC server running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to serve gRPC server: %v", err))
	}
}
