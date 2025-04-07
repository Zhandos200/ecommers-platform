package main

import (
	"order-service/infrastructure/db"
	"order-service/internal/handler"
	"order-service/internal/repository"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	database := db.NewPostgres()
	orderRepo := &repository.OrderRepository{DB: database}
	orderUsecase := &usecase.OrderUsecase{Repo: orderRepo}
	orderHandler := &handler.OrderHandler{Usecase: orderUsecase}
	orderHandler.RegisterRoutes(r)

	r.Run(":8082")
}
