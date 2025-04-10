package main

import (
	"log"

	"user-service/infrastructure"
	"user-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	db := infrastructure.NewPostgres()

	r := gin.Default()

	// Routes
	handler.RegisterRoutes(r, db)

	log.Println("✅ User service running on port 8083")
	if err := r.Run(":8083"); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
