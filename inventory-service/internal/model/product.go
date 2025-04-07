package model

type Product struct {
	ID       int     `db:"id" json:"id"`
	Name     string  `db:"name" json:"name"`
	Category string  `db:"category" json:"category"`
	Stock    int     `db:"stock" json:"stock"`
	Price    float64 `db:"price" json:"price"`
}
