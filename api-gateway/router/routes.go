package router

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	inventoryURL, _ := url.Parse("http://inventory-service:8081")
	orderURL, _ := url.Parse("http://order-service:8082")
	userURL, _ := url.Parse("http://user-service:8083")

	inventoryProxy := httputil.NewSingleHostReverseProxy(inventoryURL)
	orderProxy := httputil.NewSingleHostReverseProxy(orderURL)
	userProxy := httputil.NewSingleHostReverseProxy(userURL)

	// Inventory Routes
	r.Any("/api/products/*path", func(c *gin.Context) {
		c.Request.URL.Path = "/products" + c.Param("path")
		inventoryProxy.ServeHTTP(c.Writer, c.Request)
	})
	r.Any("/api/products", func(c *gin.Context) {
		c.Request.URL.Path = "/products"
		inventoryProxy.ServeHTTP(c.Writer, c.Request)
	})

	// Order Routes
	r.Any("/api/orders/*path", func(c *gin.Context) {
		c.Request.URL.Path = "/orders" + c.Param("path")
		orderProxy.ServeHTTP(c.Writer, c.Request)
	})
	r.Any("/api/orders", func(c *gin.Context) {
		c.Request.URL.Path = "/orders"
		orderProxy.ServeHTTP(c.Writer, c.Request)
	})

	// ✅ Handle base route (e.g., /api/users)
	r.Any("/api/users", func(c *gin.Context) {
		c.Request.URL.Path = "/" // must be single slash!
		userProxy.ServeHTTP(c.Writer, c.Request)
	})

	// ✅ Handle all subroutes like /api/users/register, /login, /profile/:id
	r.Any("/api/users/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path") // NOT "/" + c.Param("path")
		userProxy.ServeHTTP(c.Writer, c.Request)
	})

}
