package service

import (
	"github.com/mrlightwood/golang-products-api/db"
	"github.com/mrlightwood/golang-products-api/model"
)

type ProductService interface {
	CreateProduct(product *model.Product) (*int, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(id int) error
	GetProduct(id int) (*model.Product, error)
	GetProducts(category *int) ([]*model.Product, error)
}

func NewProductService(store db.Store) ProductService {
	return &ProductServiceContext{store: store}
}

type ProductServiceContext struct {
	store db.Store
}

func (psc *ProductServiceContext) GetProduct(id int) (*model.Product, error) {
	return psc.store.GetProduct(nil, id)
}

func (psc *ProductServiceContext) GetProducts(category *int) ([]*model.Product, error) {
	return psc.store.GetProducts(nil, category)
}

func (psc *ProductServiceContext) CreateProduct(product *model.Product) (*int, error) {
	tx, err := psc.store.Begin()
	if err != nil {
		return nil, err
	}
	cat, err := psc.store.CreateProduct(tx, product)
	if err != nil {
		psc.store.Rollback(tx)
		return nil, err
	}
	if err = psc.store.Commit(tx); err != nil {
		return nil, err
	}
	return cat, nil
}

func (psc *ProductServiceContext) UpdateProduct(product *model.Product) error {
	tx, err := psc.store.Begin()
	if err != nil {
		return err
	}
	err = psc.store.UpdateProduct(tx, product)
	if err != nil {
		psc.store.Rollback(tx)
		return err
	}
	if err = psc.store.Commit(tx); err != nil {
		return err
	}
	return nil
}

func (psc *ProductServiceContext) DeleteProduct(id int) error {
	tx, err := psc.store.Begin()
	if err != nil {
		return err
	}
	err = psc.store.DeleteProduct(tx, id)
	if err != nil {
		psc.store.Rollback(tx)
		return err
	}
	if err = psc.store.Commit(tx); err != nil {
		return err
	}
	return nil
}
