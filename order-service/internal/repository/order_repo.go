package repository

import (
	"order-service/internal/model"
	"time"

	"fmt"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	DB *sqlx.DB
}

func (r *OrderRepository) Create(order *model.Order) error {
	tx := r.DB.MustBegin()
	now := time.Now()

	fmt.Println("🚀 Starting transaction to insert order and items...")

	err := tx.QueryRowx(
		"INSERT INTO orders (user_id, status, created_at) VALUES ($1, $2, $3) RETURNING id",
		order.UserID, order.Status, now).Scan(&order.ID)
	if err != nil {
		tx.Rollback()
		fmt.Println("❌ Failed to insert order. Rolling back:", err)
		return err
	}
	fmt.Println("✅ Order inserted with ID:", order.ID)

	for _, item := range order.Items {
		_, err := tx.Exec("INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)",
			order.ID, item.ProductID, item.Quantity)
		if err != nil {
			tx.Rollback()
			fmt.Println("❌ Failed to insert order item. Rolling back:", err)
			return err
		}
		fmt.Printf("✅ Inserted item: ProductID=%d, Quantity=%d\n", item.ProductID, item.Quantity)
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("❌ Failed to commit transaction:", err)
		return err
	}

	fmt.Println("🎉 Transaction committed successfully.")
	return nil
}

func (r *OrderRepository) ListAll() ([]model.Order, error) {
	var orders []model.Order
	err := r.DB.Select(&orders, "SELECT * FROM orders")
	if err != nil {
		return nil, err
	}

	for i := range orders {
		r.DB.Select(&orders[i].Items, "SELECT product_id, quantity FROM order_items WHERE order_id=$1", orders[i].ID)
	}

	return orders, nil
}

func (r *OrderRepository) GetByID(id int) (*model.Order, error) {
	var order model.Order
	err := r.DB.Get(&order, "SELECT * FROM orders WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	err = r.DB.Select(&order.Items, "SELECT product_id, quantity FROM order_items WHERE order_id=$1", id)
	return &order, err
}

func (r *OrderRepository) UpdateStatus(id int, status string) error {
	_, err := r.DB.Exec("UPDATE orders SET status=$1 WHERE id=$2", status, id)
	return err
}

func (r *OrderRepository) ListByUser(userID int) ([]model.Order, error) {
	var orders []model.Order
	err := r.DB.Select(&orders, "SELECT * FROM orders WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		r.DB.Select(&orders[i].Items, "SELECT product_id, quantity FROM order_items WHERE order_id=$1", orders[i].ID)
	}

	return orders, nil
}
