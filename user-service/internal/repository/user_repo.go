package repository

import (
	"user-service/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (model.User, error)
	GetUserByID(id int) (model.User, error)
	UpdateUser(user model.User) error
	DeleteUser(id int) error
	CreatePendingUser(user *model.User, token string) error
	GetPendingByToken(token string) (*model.User, error)
	DeletePending(token string) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *model.User) error {
	return r.db.QueryRowx(
		`INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id`,
		user.Email, user.Password, user.Name,
	).Scan(&user.ID)
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

func (r *userRepo) UpdateUser(user model.User) error {
	_, err := r.db.Exec(`UPDATE users SET name=$1, email=$2 WHERE id=$3`, user.Name, user.Email, user.ID)
	return err
}

func (r *userRepo) DeleteUser(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}
func (r *userRepo) CreatePendingUser(u *model.User, token string) error {
	_, err := r.db.Exec(
		`INSERT INTO pending_users (email, name, password, token)
       VALUES ($1,$2,$3,$4)`,
		u.Email, u.Name, u.Password, token,
	)
	return err
}

func (r *userRepo) GetPendingByToken(token string) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u,
		`SELECT email, name, password FROM pending_users WHERE token=$1`,
		token,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) DeletePending(token string) error {
	_, err := r.db.Exec(`DELETE FROM pending_users WHERE token=$1`, token)
	return err
}
