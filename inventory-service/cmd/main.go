package main

import (
	"inventory-service/infrastructure/db"
	"inventory-service/internal/handler"
	"inventory-service/internal/repository"
	"inventory-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	database := db.NewPostgres()
	productRepo := &repository.ProductRepository{DB: database}
	productUsecase := &usecase.ProductUsecase{Repo: productRepo}
	productHandler := &handler.ProductHandler{Usecase: productUsecase}
	productHandler.RegisterRoutes(r)

	r.Run(":8081")
}
