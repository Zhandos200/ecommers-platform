package main

import (
	"api-gateway/cache"
	"api-gateway/logger"
	"api-gateway/middleware"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	pbInventory "api-gateway/pb/inventory"
	pbOrder "api-gateway/pb/order"
	pbUser "api-gateway/pb/user"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
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

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	logger.InitLogger()
	logger.Log.Info("🚀 API Gateway started")

	redisClient := cache.NewRedisClient()

	r := gin.Default()
	r.Static("/static", "./static")
	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"mul": func(a, b int) int {
			return a * b
		},
	}

	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
	r.SetHTMLTemplate(tmpl)

	r.Use(middleware.RequestLogger())

	invConn, _ := grpc.Dial("inventory-service:50053", grpc.WithInsecure())
	defer invConn.Close()
	inventoryClient := pbInventory.NewInventoryServiceClient(invConn)

	ordConn, _ := grpc.Dial("order-service:50052", grpc.WithInsecure())
	defer ordConn.Close()
	orderClient := pbOrder.NewOrderServiceClient(ordConn)

	usrConn, _ := grpc.Dial("user-service:50051", grpc.WithInsecure())
	defer usrConn.Close()
	userClient := pbUser.NewUserServiceClient(usrConn)

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/products", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		offset := (page - 1) * limit

		// Call gRPC to get all products
		res, err := inventoryClient.ListProducts(c, &pbInventory.Empty{})
		if err != nil {
			c.String(500, "Error loading products: %v", err)
			return
		}

		// Convert to local struct
		var allProducts []Product
		for _, p := range res.Products {
			allProducts = append(allProducts, Product{
				ID:       int(p.Id),
				Name:     p.Name,
				Category: p.Category,
				Stock:    int(p.Stock),
				Price:    float64(p.Price),
			})
		}

		// Apply pagination
		total := len(allProducts)
		end := offset + limit
		if end > total {
			end = total
		}
		if offset > total {
			offset = total
		}
		products := allProducts[offset:end]

		c.HTML(http.StatusOK, "products.html", gin.H{
			"Products": products,
			"Page":     page,
			"Limit":    limit,
			"Total":    total,
		})
	})

	r.POST("/products", func(c *gin.Context) {
		var req Product
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
			return
		}

		grpcReq := &pbInventory.Product{
			Name:     req.Name,
			Category: req.Category,
			Stock:    int32(req.Stock),
			Price:    float32(req.Price),
		}

		grpcRes, err := inventoryClient.CreateProduct(context.Background(), grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
			return
		}

		response := Product{
			ID:       int(grpcRes.Id),
			Name:     grpcRes.Name,
			Category: grpcRes.Category,
			Stock:    int(grpcRes.Stock),
			Price:    float64(grpcRes.Price),
		}
		c.JSON(http.StatusCreated, response)
	})

	r.GET("/products/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid product ID"})
			return
		}

		res, err := inventoryClient.GetProduct(context.Background(), &pbInventory.ProductID{Id: int64(id)})
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get product", "details": err.Error()})
			return
		}

		c.JSON(200, Product{
			ID:       int(res.Id),
			Name:     res.Name,
			Category: res.Category,
			Stock:    int(res.Stock),
			Price:    float64(res.Price),
		})
	})

	r.PUT("/products/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid product ID"})
			return
		}

		var input Product
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid product data"})
			return
		}

		req := &pbInventory.Product{
			Id:       int64(id),
			Name:     input.Name,
			Category: input.Category,
			Stock:    int32(input.Stock),
			Price:    float32(input.Price),
		}

		_, err = inventoryClient.UpdateProduct(context.Background(), req)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to update product", "details": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Product updated"})
	})

	r.DELETE("/products/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid product ID"})
			return
		}

		_, err = inventoryClient.DeleteProduct(context.Background(), &pbInventory.ProductID{Id: int64(id)})
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete product", "details": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Product deleted"})
	})

	r.GET("/orders", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		offset := (page - 1) * limit

		res, err := orderClient.ListOrders(c, &pbOrder.UserOrdersRequest{})
		if err != nil {
			c.String(500, "Error loading orders: %v", err)
			return
		}

		var allOrders []Order
		for _, o := range res.Orders {
			var items []OrderItem
			for _, item := range o.Items {
				items = append(items, OrderItem{
					ProductID: int(item.ProductId),
					Quantity:  int(item.Quantity),
				})
			}
			allOrders = append(allOrders, Order{
				ID:        int(o.Id),
				UserID:    fmt.Sprint(o.UserId),
				Status:    o.Status,
				CreatedAt: o.CreatedAt,
				Items:     items,
			})
		}

		// Apply pagination
		total := len(allOrders)
		end := offset + limit
		if end > total {
			end = total
		}
		if offset > total {
			offset = total
		}
		orders := allOrders[offset:end]

		c.HTML(http.StatusOK, "orders.html", gin.H{
			"Orders": orders,
			"Page":   page,
			"Limit":  limit,
			"Total":  total,
		})
	})

	r.POST("/orders", func(c *gin.Context) {
		var input struct {
			UserID string `json:"user_id"`
			Items  []struct {
				ProductID int `json:"product_id"`
				Quantity  int `json:"quantity"`
			} `json:"items"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order format"})
			return
		}

		userIDInt, err := strconv.ParseInt(input.UserID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var items []*pbOrder.OrderItem
		for _, item := range input.Items {
			items = append(items, &pbOrder.OrderItem{
				ProductId: int64(item.ProductID),
				Quantity:  int32(item.Quantity),
			})
		}

		req := &pbOrder.OrderRequest{
			UserId: userIDInt,
			Items:  items,
		}

		res, err := orderClient.CreateOrder(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":         res.Id,
			"user_id":    res.UserId,
			"status":     res.Status,
			"created_at": res.CreatedAt,
			"items":      res.Items,
		})
	})

	r.GET("/orders/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid order ID"})
			return
		}

		res, err := orderClient.GetOrder(context.Background(), &pbOrder.OrderID{Id: int64(id)})
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get order", "details": err.Error()})
			return
		}

		var items []OrderItem
		for _, item := range res.Items {
			items = append(items, OrderItem{
				ProductID: int(item.ProductId),
				Quantity:  int(item.Quantity),
			})
		}

		c.JSON(200, Order{
			ID:        int(res.Id),
			UserID:    fmt.Sprint(res.UserId),
			Status:    res.Status,
			CreatedAt: res.CreatedAt,
			Items:     items,
		})
	})

	r.PATCH("/orders/:id/status", func(c *gin.Context) {
		var input struct {
			Status string `json:"status"`
		}
		if err := c.ShouldBindJSON(&input); err != nil || input.Status == "" {
			c.JSON(400, gin.H{"error": "Missing or invalid status"})
			return
		}

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid order ID"})
			return
		}

		_, err = orderClient.UpdateOrderStatus(context.Background(), &pbOrder.StatusUpdate{
			Id:     int64(id),
			Status: input.Status,
		})
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to update order status", "details": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Order status updated"})
	})

	r.GET("/users", func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			c.String(400, "Missing user ID")
			return
		}
		c.Redirect(http.StatusFound, "/users/"+id)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		cacheKey := "user:" + id
		fmt.Println("Looking for key in Redis:", cacheKey)

		// Try to get user from Redis
		cached, err := cache.GetCache(redisClient, cacheKey)
		if err == nil {
			fmt.Println("Cache HIT for:", cacheKey)
			var user User
			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				fmt.Println("Unmarshalled cached user:", user)
				c.HTML(200, "users.html", gin.H{"User": user})
				return
			} else {
				fmt.Println("Failed to unmarshal cached user:", err)
			}
		} else {
			fmt.Println("Cache MISS for:", cacheKey, "→", err)
		}

		// If not in cache, call gRPC
		fmt.Println("📡 Fetching user from gRPC service...")
		u64, _ := strconv.ParseInt(id, 10, 64)
		res, err := userClient.GetUserProfile(c, &pbUser.UserID{Id: u64})
		if err != nil {
			fmt.Printf("gRPC error for user %s: %v\n", id, err)
			c.String(500, "Error loading user: %v", err)
			return
		}

		user := User{ID: int(res.Id), Email: res.Email, Name: res.Name}
		fmt.Println("gRPC returned user:", user)

		// Marshal and store in Redis
		userJson, _ := json.Marshal(user)
		err = cache.SetCache(redisClient, cacheKey, string(userJson), 5*time.Minute)
		if err != nil {
			fmt.Println("Failed to set cache:", err)
		} else {
			fmt.Println("📝 Cached user to Redis with key:", cacheKey)
		}

		c.HTML(200, "users.html", gin.H{"User": user})
	})

	r.GET("/users/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})

	r.POST("/users/register", func(c *gin.Context) {
		var input struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		req := &pbUser.UserRequest{
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password,
		}

		res, err := userClient.RegisterUser(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":    res.Id,
			"name":  res.Name,
			"email": res.Email,
		})
	})

	r.GET("/users/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	r.POST("/users/login", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")

		_, err := userClient.AuthenticateUser(c, &pbUser.AuthRequest{
			Email:    email,
			Password: password,
		})
		if err != nil {
			c.String(401, "Invalid credentials")
			return
		}
		c.Redirect(302, "/")
	})

	r.PATCH("/users/:id", func(c *gin.Context) {
		var input struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID"})
			return
		}

		req := &pbUser.UserRequest{
			Id:       int64(id),
			Email:    input.Email,
			Name:     input.Name,
			Password: input.Password,
		}

		res, err := userClient.UpdateUser(context.Background(), req)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to update user", "details": err.Error()})
			return
		}

		// Invalidate cache
		cache.DeleteCache(redisClient, "user:"+idParam)

		c.JSON(200, gin.H{
			"id":    res.Id,
			"email": res.Email,
			"name":  res.Name,
		})
	})

	r.DELETE("/users/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID"})
			return
		}

		_, err = userClient.DeleteUser(context.Background(), &pbUser.UserID{Id: int64(id)})
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete user", "details": err.Error()})
			return
		}

		// Invalidate cache
		cache.DeleteCache(redisClient, "user:"+idParam)

		c.JSON(200, gin.H{"message": "User deleted"})
	})

	r.Run(":8080")
}
