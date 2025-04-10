package router

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	inventoryURL, _ := url.Parse("http://inventory-service:8081")
	orderURL, _ := url.Parse("http://order-service:8082")

	inventoryProxy := httputil.NewSingleHostReverseProxy(inventoryURL)
	orderProxy := httputil.NewSingleHostReverseProxy(orderURL)

	r.Any("/api/products", func(c *gin.Context) {
		c.Request.URL.Path = "/products"
		inventoryProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Any("/api/products/*any", func(c *gin.Context) {
		c.Request.URL.Path = "/products" + c.Param("any")
		inventoryProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Any("/api/orders", func(c *gin.Context) {
		c.Request.URL.Path = "/orders"
		orderProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Any("/api/orders/*any", func(c *gin.Context) {
		c.Request.URL.Path = "/orders" + c.Param("any")
		orderProxy.ServeHTTP(c.Writer, c.Request)
	})
}
