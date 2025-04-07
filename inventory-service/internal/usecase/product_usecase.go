package usecase

import (
	"inventory-service/internal/model"
)

type ProductRepo interface {
	Create(product *model.Product) error
	GetByID(id int) (*model.Product, error)
	Update(id int, product *model.Product) error
	Delete(id int) error
	List(category string, limit, offset int) ([]model.Product, error)
}

type ProductUsecase struct {
	Repo ProductRepo
}

func (u *ProductUsecase) Create(p *model.Product) error {
	return u.Repo.Create(p)
}

func (u *ProductUsecase) GetByID(id int) (*model.Product, error) {
	return u.Repo.GetByID(id)
}

func (u *ProductUsecase) Update(id int, p *model.Product) error {
	return u.Repo.Update(id, p)
}

func (u *ProductUsecase) Delete(id int) error {
	return u.Repo.Delete(id)
}

func (u *ProductUsecase) List(category string, limit, offset int) ([]model.Product, error) {
	return u.Repo.List(category, limit, offset)
}
