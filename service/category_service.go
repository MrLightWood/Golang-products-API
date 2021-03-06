package service

import (
	"github.com/mrlightwood/golang-products-api/db"
	"github.com/mrlightwood/golang-products-api/model"
)

type CategoryService interface {
	CreateCategory(category *model.Category) (*int, error)
	UpdateCategory(category *model.Category) error
	DeleteCategory(id int) error
	GetCategory(id int) (*model.Category, error)
	GetCategories() ([]*model.Category, error)
}

type CategoryServiceContext struct {
	store db.Store
}

func NewCategoryService(store db.Store) CategoryService {
	return &CategoryServiceContext{store: store}
}

func (csc *CategoryServiceContext) GetCategory(id int) (*model.Category, error) {
	return csc.store.GetCategory(nil, id)
}

func (csc *CategoryServiceContext) GetCategories() ([]*model.Category, error) {
	return csc.store.GetCategories(nil)
}

func (csc *CategoryServiceContext) CreateCategory(category *model.Category) (*int, error) {
	tx, err := csc.store.Begin()
	if err != nil {
		return nil, err
	}
	cat, err := csc.store.CreateCategory(tx, category)
	if err != nil {
		csc.store.Rollback(tx)
		return nil, err
	}
	if err = csc.store.Commit(tx); err != nil {
		return nil, err
	}
	return cat, nil
}

func (csc *CategoryServiceContext) UpdateCategory(category *model.Category) error {
	tx, err := csc.store.Begin()
	if err != nil {
		return err
	}
	err = csc.store.UpdateCategory(tx, category)
	if err != nil {
		csc.store.Rollback(tx)
		return err
	}
	if err = csc.store.Commit(tx); err != nil {
		return err
	}
	return nil
}

func (csc *CategoryServiceContext) DeleteCategory(id int) error {
	tx, err := csc.store.Begin()
	if err != nil {
		return err
	}
	err = csc.store.DeleteCategory(tx, id)
	if err != nil {
		csc.store.Rollback(tx)
		return err
	}
	if err = csc.store.Commit(tx); err != nil {
		return err
	}
	return nil
}
