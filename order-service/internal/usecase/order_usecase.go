package usecase

import "order-service/internal/model"

type OrderRepo interface {
	Create(order *model.Order) error
	GetByID(id int) (*model.Order, error)
	UpdateStatus(id int, status string) error
	ListByUser(userID string) ([]model.Order, error)
	ListAll() ([]model.Order, error)
}

type OrderUsecase struct {
	Repo OrderRepo
}

func (u *OrderUsecase) ListAll() ([]model.Order, error) {
	return u.Repo.ListAll()
}

func (u *OrderUsecase) Create(o *model.Order) error {
	return u.Repo.Create(o)
}

func (u *OrderUsecase) GetByID(id int) (*model.Order, error) {
	return u.Repo.GetByID(id)
}

func (u *OrderUsecase) UpdateStatus(id int, status string) error {
	return u.Repo.UpdateStatus(id, status)
}

func (u *OrderUsecase) ListByUser(userID string) ([]model.Order, error) {
	return u.Repo.ListByUser(userID)
}
