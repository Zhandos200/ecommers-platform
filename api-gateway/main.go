package main

import (
	"api-gateway/middleware"
	"api-gateway/router"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Stock    int     `json:"stock"`
	Price    float64 `json:"price"`
}

type OrderItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Order struct {
	ID        int         `json:"id"`
	UserID    string      `json:"user_id"`
	Status    string      `json:"status"`
	CreatedAt string      `json:"created_at"`
	Items     []OrderItem `json:"items"`
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Use(middleware.RequestLogger())

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/products", func(c *gin.Context) {
		resp, err := http.Get("http://localhost:8080/api/products") // use API Gateway
		if err != nil {
			c.String(500, "Error loading products")
			return
		}
		defer resp.Body.Close()

		var products []Product
		json.NewDecoder(resp.Body).Decode(&products)

		c.HTML(200, "products.html", gin.H{"Products": products})
	})

	r.GET("/orders", func(c *gin.Context) {
		resp, err := http.Get("http://localhost:8080/api/orders")
		if err != nil {
			c.String(500, "Error loading orders")
			return
		}
		defer resp.Body.Close()

		var orders []Order
		json.NewDecoder(resp.Body).Decode(&orders)

		c.HTML(200, "orders.html", gin.H{"Orders": orders})
	})

	router.SetupRoutes(r)

	r.Run(":8080") // API Gateway listens on port 8080
}
