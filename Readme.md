# E-Commerce Platform (gRPC Microservices + REST API Gateway)

This project is a simplified e-commerce platform built using **Go**, **gRPC**, and **Gin**. It features a **microservices architecture** with internal gRPC communication and a RESTful API gateway.

---

## Architecture


---

## Services

### 1. Inventory Service
- gRPC methods: `CreateProduct`, `GetProduct`, `UpdateProduct`, `DeleteProduct`, `ListProducts`
- PostgreSQL for product data

### 2. Order Service
- gRPC methods: `CreateOrder`, `ListOrders`
- Calls **Inventory Service** internally to check stock before creating an order

### 3. User Service
- gRPC methods: `RegisterUser`, `AuthenticateUser`, `GetUserProfile`, `UpdateUser`, `DeleteUser`
- Secure password hashing with bcrypt

---

## API Gateway
- Exposes RESTful endpoints:
  - `POST /products`
  - `GET /products`
  - `POST /orders`
  - `GET /orders`
  - `POST /users/register`
  - `POST /users/login`
  - `PATCH /users/:id`
  - `DELETE /users/:id`

---

## Technologies Used

- Go 1.22+
- gRPC & Protocol Buffers
- Gin Web Framework
- PostgreSQL
- Docker & Docker Compose

---

## Setup & Run

### 1. Clone the repository:
```bash
git clone https://github.com/your-username/ecommers-platform.git
cd ecommers-platform


2. Start all services:
docker-compose up --build

3. Access API Gateway:
Open browser: http://localhost:8080

Create Product (POST /products)
{
  "name": "MacBook Pro",
  "category": "Laptops",
  "stock": 50,
  "price": 1999.99
}


Create Order (POST /orders)
{
  "user_id": "1",
  "items": [
    {"product_id": 1, "quantity": 2 }
  ]
}

Notes
Make sure PostgreSQL is running and accessible.

Proto files are shared across services.

For local dev, consider using replace in go.mod or copying proto-generated files.

