package main

import (
	"api-gateway/middleware"
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
		email := c.PostForm("email")
		name := c.PostForm("name")
		password := c.PostForm("password")

		_, err := userClient.RegisterUser(c, &pbUser.UserRequest{
			Email:    email,
			Name:     name,
			Password: password,
		})
		if err != nil {
			c.String(500, "Failed to register user: %v", err)
			return
		}
		c.Redirect(302, "/")
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

	r.Run(":8080")
}
