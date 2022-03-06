package test

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mrlightwood/golang-products-api/config"
	"github.com/mrlightwood/golang-products-api/db"
	"github.com/mrlightwood/golang-products-api/model"
	"github.com/stretchr/testify/assert"
)

var st db.Store

func init() {
	conf, _ := config.NewConfig("../config/config_test.yaml")
	st, _ = db.NewStore(conf)
}

func TestStore_CreateCategory(t *testing.T) {
	tx, _ := st.Begin()
	defer st.Rollback(tx)
	id, err := st.CreateCategory(tx, &model.Category{Name: "test"})
	assert.NoError(t, err)
	assert.NotNil(t, id)
}

func TestStore_GetCategory(t *testing.T) {
	c, err := st.GetCategory(nil, -1)
	assert.Nil(t, err)
	assert.Nil(t, c)
	tx, _ := st.Begin()
	defer st.Rollback(tx)
	id, _ := st.CreateCategory(tx, &model.Category{Name: "test"})
	cat, _ := st.GetCategory(tx, *id)
	assert.Equal(t, *id, cat.Id)
}

func TestStore_GetCategories(t *testing.T) {
	_, err := st.GetCategories(nil)
	assert.NoError(t, err)
	tx, _ := st.Begin()
	defer st.Rollback(tx)
	st.CreateCategory(tx, &model.Category{Name: "test"})
	res, _ := st.GetCategories(tx)
	assert.NotEmpty(t, res)
}

func TestStore_UpdateCategory(t *testing.T) {
	tx, _ := st.Begin()
	defer st.Rollback(tx)
	id, _ := st.CreateCategory(tx, &model.Category{Name: "test"})
	cat, _ := st.GetCategory(tx, *id)
	cat.Name = "test2"
	st.UpdateCategory(tx, cat)
	cat2, _ := st.GetCategory(tx, cat.Id)
	assert.Equal(t, cat2.Name, "test2")
}

func TestStore_DeleteCategory(t *testing.T) {
	tx, _ := st.Begin()
	defer st.Rollback(tx)
	id, _ := st.CreateCategory(tx, &model.Category{Name: "test"})
	st.DeleteCategory(tx, *id)
	cat2, _ := st.GetCategory(tx, *id)
	assert.Nil(t, cat2)
}

func TestStore_CreateProduct(t *testing.T) {
	tx, _ := st.Begin()
	defer st.Rollback(tx)
	category, _ := st.CreateCategory(tx, &model.Category{Name: "test"})
	p, err := st.CreateProduct(tx, &model.Product{Name: "test_name", Description: "test_description", Category: *category, Price: 102.5})
	assert.NoError(t, err)
	product, err := st.GetProduct(tx, *p)
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, product.Name, "test_name")
	assert.Equal(t, product.Description, "test_description")
	assert.Equal(t, product.Category, *category)
	assert.Equal(t, product.Price, 102.5)
}

func TestStore_GetProduct(t *testing.T) {
	tx, _ := st.Begin()
	defer st.Rollback(tx)
	category, _ := st.CreateCategory(tx, &model.Category{Name: "test"})
	p, _ := st.CreateProduct(tx, &model.Product{Name: "test_name", Description: "test_description", Category: *category, Price: 102.5})
	ps, err := st.GetProduct(tx, *p)
	assert.Nil(t, err)
	assert.NotNil(t, ps)
	ps, err = st.GetProduct(tx, -1)
	assert.Nil(t, err)
	assert.Nil(t, ps)
}

func TestStore_GetProducts(t *testing.T) {
	tx, _ := st.Begin()
	defer st.Rollback(tx)
	category, _ := st.CreateCategory(tx, &model.Category{Name: "test"})
	st.CreateProduct(tx, &model.Product{Name: "test_name", Description: "test_description", Category: *category, Price: 102.5})
	st.CreateProduct(tx, &model.Product{Name: "test_name2", Description: "test_description2", Category: *category, Price: 102.52})
	ps, err := st.GetProducts(tx, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, ps)
	ps, _ = st.GetProducts(tx, category)
	assert.Len(t, ps, 2)
}

func TestStore_UpdateProduct(t *testing.T) {
	tx, _ := st.Begin()
	defer st.Rollback(tx)
	category, _ := st.CreateCategory(tx, &model.Category{Name: "test"})
	category2, _ := st.CreateCategory(tx, &model.Category{Name: "test"})
	id, _ := st.CreateProduct(tx, &model.Product{Name: "test_name", Description: "test_description", Category: *category, Price: 102.5})
	p, _ := st.GetProduct(tx, *id)
	p.Name = "test_name2"
	p.Description = "test_description2"
	p.Category = *category2
	p.Price = 102.7
	err := st.UpdateProduct(tx, p)
	assert.Nil(t, err)
	p2, _ := st.GetProduct(tx, p.Id)
	assert.Equal(t, p2.Name, "test_name2")
	assert.Equal(t, p2.Description, "test_description2")
	assert.Equal(t, p2.Category, *category2)
	assert.Equal(t, p2.Price, 102.7)
}

func TestStore_DeleteProduct(t *testing.T) {
	tx, _ := st.Begin()
	defer st.Rollback(tx)
	id, _ := st.CreateCategory(tx, &model.Category{Name: "test"})
	product, _ := st.CreateProduct(tx, &model.Product{Name: "test_name", Description: "test_description", Category: *id, Price: 102.5})
	err := st.DeleteProduct(tx, *product)
	assert.Nil(t, err)
	p, _ := st.GetProduct(tx, *product)
	assert.Nil(t, p)
}
