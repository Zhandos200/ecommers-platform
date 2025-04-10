package repository

import (
	"user-service/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user model.User) error
	GetUserByEmail(email string) (model.User, error)
	GetUserByID(id int) (model.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user model.User) error {
	_, err := r.db.Exec(`INSERT INTO users (email, password, name) VALUES ($1, $2, $3)`, user.Email, user.Password, user.Name)
	return err
}

func (r *userRepo) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Get(&user, `SELECT * FROM users WHERE email=$1`, email)
	return user, err
}

func (r *userRepo) GetUserByID(id int) (model.User, error) {
	var user model.User
	err := r.db.Get(&user, `SELECT id, email, name FROM users WHERE id=$1`, id)
	return user, err
}
