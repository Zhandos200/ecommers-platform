package main

import (
	"api-gateway/middleware"
	"context"
	"fmt"
	"net/http"
	"strconv"

	pbInventory "api-gateway/pb/inventory"
	pbOrder "api-gateway/pb/order"
	pbUser "api-gateway/pb/user"

	"github.com/gin-gonic/gin"
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
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
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
		res, err := inventoryClient.ListProducts(c, &pbInventory.Empty{})
		if err != nil {
			c.String(500, "Error loading products: %v", err)
			return
		}

		var products []Product
		for _, p := range res.Products {
			products = append(products, Product{
				ID:       int(p.Id),
				Name:     p.Name,
				Category: p.Category,
				Stock:    int(p.Stock),
				Price:    float64(p.Price),
			})
		}

		c.HTML(200, "products.html", gin.H{"Products": products})
	})

	r.GET("/orders", func(c *gin.Context) {
		res, err := orderClient.ListOrders(c, &pbOrder.UserOrdersRequest{})
		if err != nil {
			c.String(500, "Error loading orders: %v", err)
			return
		}

		var orders []Order
		for _, o := range res.Orders {
			var items []OrderItem
			for _, item := range o.Items {
				items = append(items, OrderItem{
					ProductID: int(item.ProductId),
					Quantity:  int(item.Quantity),
				})
			}
			orders = append(orders, Order{
				ID:        int(o.Id),
				UserID:    fmt.Sprint(o.UserId),
				Status:    o.Status,
				CreatedAt: o.CreatedAt,
				Items:     items,
			})
		}

		c.HTML(200, "orders.html", gin.H{"Orders": orders})
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
		u64, _ := strconv.ParseInt(id, 10, 64)
		res, err := userClient.GetUserProfile(c, &pbUser.UserID{Id: u64})
		if err != nil {
			c.String(500, "Error loading user: %v", err)
			return
		}
		c.HTML(200, "users.html", gin.H{"User": User{ID: int(res.Id), Email: res.Email, Name: res.Name}})
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

		// Convert gRPC response back to REST format
		response := Product{
			ID:       int(grpcRes.Id),
			Name:     grpcRes.Name,
			Category: grpcRes.Category,
			Stock:    int(grpcRes.Stock),
			Price:    float64(grpcRes.Price),
		}

		c.JSON(http.StatusCreated, response)
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

		c.JSON(200, gin.H{"message": "User deleted"})
	})

	r.Run(":8080")
}
