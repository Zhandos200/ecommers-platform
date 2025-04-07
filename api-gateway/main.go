package main

import (
	"api-gateway/middleware"
	"api-gateway/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middleware.RequestLogger())

	router.SetupRoutes(r)

	r.Run(":8080") // API Gateway listens on port 8080
}
