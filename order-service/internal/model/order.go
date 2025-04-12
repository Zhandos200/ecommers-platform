package model

import "time"

type Order struct {
	ID        int         `db:"id" json:"id"`
	UserID    int         `db:"user_id" json:"user_id"` // ✅ change from string to int
	Status    string      `db:"status" json:"status"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	Items     []OrderItem `json:"items"`
}

type OrderItem struct {
	OrderID   int `db:"order_id" json:"-"`
	ProductID int `db:"product_id" json:"product_id"`
	Quantity  int `db:"quantity" json:"quantity"`
}
