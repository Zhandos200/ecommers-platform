# 🛒 E-Commerce Platform (gRPC Microservices + REST API Gateway + Observability)

A modern microservice-based e-commerce platform built using **Go**, **gRPC**, **Gin**, **Docker**, and **Prometheus/Grafana/Tempo/Loki** stack. It supports **user management**, **product inventory**, **order processing**, **RESTful UI**, **Redis caching**, and **observability** with logs, metrics, and traces.

---

## 🧱 Architecture

      [ HTML + CSS Pages ]
             |
    ┌────────▼────────┐
    │   API Gateway   │ - REST (Gin)
    └──┬────┬────┬────┘
       │    │    │
 gRPC  │    │    │  gRPC
   ┌───▼┐ ┌─▼───┐ ┌──▼───┐
   │User│ │Order│ │Inventory│
   └────┘ └─────┘ └───────┘

Redis for caching
PostgreSQL for persistence
NATS for async events
Tempo, Loki, Prometheus, Grafana for observability


---

## 📦 Microservices

### 1. **User Service**
- Register/Login (secure bcrypt password)
- gRPC methods:
  - `RegisterUser`, `AuthenticateUser`
  - `GetUserProfile`, `UpdateUser`, `DeleteUser`

### 2. **Inventory Service**
- Full Product CRUD
- gRPC methods:
  - `CreateProduct`, `GetProduct`, `UpdateProduct`, `DeleteProduct`, `ListProducts`

### 3. **Order Service**
- Create/List Orders
- Calls inventory service internally to update stock
- gRPC methods:
  - `CreateOrder`, `GetOrder`, `ListOrders`, `UpdateOrderStatus`

---

## 🌐 REST API Gateway (Gin)

Serves HTML pages and REST endpoints.

### HTML Pages
- `/` – Home
- `/products` – Product listing with pagination
- `/orders` – Order listing with pagination
- `/users/:id` – User profile (with Redis caching)
- `/users/register`, `/users/login` – Forms

### Product Endpoints
- `GET /products`
- `GET /products/:id`
- `POST /products`
- `PUT /products/:id`
- `DELETE /products/:id`

### Order Endpoints
- `GET /orders`
- `GET /orders/:id`
- `POST /orders`
- `PATCH /orders/:id/status`

### User Endpoints
- `GET /users?id=1` → redirects to `/users/1`
- `GET /users/:id`
- `POST /users/register`
- `POST /users/login`
- `PATCH /users/:id`
- `DELETE /users/:id`

---

## 📈 Observability & Monitoring

| Tool       | Purpose             | Access                             |
|------------|---------------------|-------------------------------------|
| **Grafana** | Unified dashboards  | http://localhost:3000               |
| **Prometheus** | Metrics collection | http://localhost:9090               |
| **Tempo**  | Distributed tracing | Integrated in Grafana               |
| **Loki**   | Log aggregation     | View logs via Grafana Explore tab   |

- `/metrics` endpoint exposed at `:2112` for Prometheus
- Gin middleware logs structured to Loki
- Tempo captures gRPC + HTTP traces

---

## 🧰 Technologies Used

- **Go 1.24**
- **Gin (REST)**
- **gRPC + Protobuf**
- **PostgreSQL**
- **Redis**
- **Docker + Compose**
- **Prometheus, Grafana, Tempo, Loki**

---

## 🏁 Getting Started

### 1. Clone the repo
```bash
git clone https://github.com/your-username/ecommers-platform.git
cd ecommers-platform

Run all services
  docker-compose up --build


 Open in browser
    Gateway: http://localhost:8080
    Grafana: http://localhost:3000

.env(prompt, add to root):
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB_USERS=users
POSTGRES_DB_ORDERS=orders
POSTGRES_DB_INVENTORY=inventory
