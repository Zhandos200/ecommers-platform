package repository

import (
	"inventory-service/internal/model"

	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	DB *sqlx.DB
}

func (r *ProductRepository) Create(p *model.Product) error {
	query := `INSERT INTO products (name, category, stock, price) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, p.Name, p.Category, p.Stock, p.Price)
	return err
}

func (r *ProductRepository) GetByID(id int) (*model.Product, error) {
	var p model.Product
	err := r.DB.Get(&p, "SELECT * FROM products WHERE id = $1", id)
	return &p, err
}

func (r *ProductRepository) Update(id int, p *model.Product) error {
	query := `UPDATE products SET name=$1, category=$2, stock=$3, price=$4 WHERE id=$5`
	_, err := r.DB.Exec(query, p.Name, p.Category, p.Stock, p.Price, id)
	return err
}

func (r *ProductRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}

func (r *ProductRepository) List(category string, limit, offset int) ([]model.Product, error) {
	var products []model.Product
	query := `SELECT * FROM products WHERE ($1 = '' OR category = $1) LIMIT $2 OFFSET $3`
	err := r.DB.Select(&products, query, category, limit, offset)
	return products, err
}
