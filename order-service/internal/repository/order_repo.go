package repository

import (
	"order-service/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	DB *sqlx.DB
}

func (r *OrderRepository) Create(order *model.Order) error {
	tx := r.DB.MustBegin()
	now := time.Now()
	err := tx.QueryRowx(
		"INSERT INTO orders (user_id, status, created_at) VALUES ($1, $2, $3) RETURNING id",
		order.UserID, order.Status, now).Scan(&order.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err := tx.Exec("INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)",
			order.ID, item.ProductID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
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

func (r *OrderRepository) ListByUser(userID string) ([]model.Order, error) {
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
