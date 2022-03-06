package db

import (
	"database/sql"
	"errors"

	"github.com/mrlightwood/golang-products-api/config"
	"github.com/mrlightwood/golang-products-api/model"
)

type Store interface {
	// Begin transaction
	Begin() (*sql.Tx, error)
	// Close storage
	Close() error
	// Commit transaction
	Commit(tx *sql.Tx) error
	// Rollback transaction
	Rollback(tx *sql.Tx) error
	// Get product by id
	GetProduct(tx *sql.Tx, id int) (*model.Product, error)
	// Get all products
	GetProducts(tx *sql.Tx, category *int) ([]*model.Product, error)
	// Create a new product
	CreateProduct(tx *sql.Tx, product *model.Product) (*int, error)
	// Update an existing product
	UpdateProduct(tx *sql.Tx, product *model.Product) error
	// Delete an existing product
	DeleteProduct(tx *sql.Tx, id int) error
	// Get category by id
	GetCategory(tx *sql.Tx, id int) (*model.Category, error)
	// Get all categories
	GetCategories(tx *sql.Tx) ([]*model.Category, error)
	// Create an existing category
	CreateCategory(tx *sql.Tx, category *model.Category) (*int, error)
	// Update an existing category
	UpdateCategory(tx *sql.Tx, category *model.Category) error
	// Delete an existing category
	DeleteCategory(tx *sql.Tx, id int) error
}

type StoreContext struct {
	db *sql.DB
}

func initialize(db *sql.DB, tx *sql.Tx) error {
	var query = `CREATE TABLE IF NOT EXISTS "category" (
		"id"	INTEGER NOT NULL,
		"name"	TEXT NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	CREATE TABLE IF NOT EXISTS "product" (
		"id"	INTEGER NOT NULL,
		"name"	TEXT NOT NULL,
		"description"	TEXT,
		"category"	INTEGER,
		"price"	REAL NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT)
	);`

	var err error

	if tx != nil {
		_, err = tx.Exec(query)
	} else {
		_, err = db.Exec(query)
	}
	if err != nil {
		return err
	}

	return nil
}

func NewStore(conf *config.Config) (Store, error) {
	db, err := sql.Open("sqlite3", conf.Store.Dbpath)

	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	err = initialize(db, nil)
	if err != nil {
		return nil, err
	}

	return &StoreContext{db}, nil
}

func (sc *StoreContext) Close() error {
	return sc.db.Close()
}

func (sc *StoreContext) Begin() (*sql.Tx, error) {
	return sc.db.Begin()
}

func (sc *StoreContext) Commit(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("transaction is nil")
	}

	return tx.Commit()
}

func (sc *StoreContext) Rollback(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("transaction is nil")
	}
	return tx.Rollback()
}

func (sc *StoreContext) GetCategory(tx *sql.Tx, id int) (*model.Category, error) {
	var query = "SELECT id, name FROM category WHERE id= $1;"
	var row *sql.Row

	if tx != nil {
		row = tx.QueryRow(query, id)
	} else {
		row = sc.db.QueryRow(query, id)
	}
	category := &model.Category{}
	if err := row.Scan(&category.Id, &category.Name); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		} else {
			return nil, nil
		}
	}
	return category, nil
}

func (sc *StoreContext) GetCategories(tx *sql.Tx) ([]*model.Category, error) {
	query := "SELECT id, name FROM category;"
	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(query)
	} else {
		rows, err = sc.db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*model.Category
	for rows.Next() {
		category := &model.Category{}
		if err := rows.Scan(&category.Id, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (sc *StoreContext) CreateCategory(tx *sql.Tx, category *model.Category) (*int, error) {
	var query = "INSERT INTO category(name) VALUES($1) RETURNING id;"
	var id int
	var err error
	if tx != nil {
		err = tx.QueryRow(query, category.Name).Scan(&id)
	} else {
		err = sc.db.QueryRow(query, category.Name).Scan(&id)
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (sc *StoreContext) UpdateCategory(tx *sql.Tx, category *model.Category) error {
	query := "UPDATE category SET name =$1 WHERE id = $2;"
	var res sql.Result
	var err error
	if tx != nil {
		res, err = tx.Exec(query, category.Name, category.Id)
	} else {
		res, err = sc.db.Exec(query, category.Name, category.Id)
	}
	if err != nil {
		return err
	}
	if a, err := res.RowsAffected(); err != nil {
		return err
	} else if a == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (sc *StoreContext) DeleteCategory(tx *sql.Tx, id int) error {
	query := "DELETE FROM category WHERE id = $1;"
	var res sql.Result
	var err error
	if tx != nil {
		res, err = tx.Exec(query, id)
	} else {
		res, err = sc.db.Exec(query, id)
	}
	if err != nil {
		return err
	}
	if a, err := res.RowsAffected(); err != nil {
		return err
	} else if a == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (sc *StoreContext) GetProduct(tx *sql.Tx, id int) (*model.Product, error) {
	var query = "SELECT id, name, description, category, price FROM product WHERE id= $1;"
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow(query, id)
	} else {
		row = sc.db.QueryRow(query, id)
	}
	product := &model.Product{}
	if err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Category, &product.Price); err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		} else {
			return nil, nil
		}
	}
	return product, nil
}

func (sc *StoreContext) GetProducts(tx *sql.Tx, category *int) ([]*model.Product, error) {
	var query string
	var rows *sql.Rows
	var err error
	if category == nil {
		query = "SELECT id, name, description, category, price FROM product;"
		if tx != nil {
			rows, err = tx.Query(query)
		} else {
			rows, err = sc.db.Query(query)
		}
	} else {
		query = "SELECT id, name, description, category, price FROM product WHERE category= $1;"
		if tx != nil {
			rows, err = tx.Query(query, category)
		} else {
			rows, err = sc.db.Query(query, category)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*model.Product

	for rows.Next() {
		product := &model.Product{}
		if err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Category, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (sc *StoreContext) CreateProduct(tx *sql.Tx, product *model.Product) (*int, error) {
	var query = "INSERT INTO product( name, description, category, price) VALUES($1, $2, $3, $4) RETURNING id;"
	var id int
	var err error
	if tx != nil {
		err = tx.QueryRow(query, product.Name, product.Description, product.Category, product.Price).Scan(&id)
	} else {
		err = sc.db.QueryRow(query, product.Name, product.Description, product.Category, product.Price).Scan(&id)
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (sc *StoreContext) UpdateProduct(tx *sql.Tx, product *model.Product) error {
	query := "UPDATE product SET name=$1, description=$2, category=$3, price=$4  WHERE id = $5;"
	var res sql.Result
	var err error
	if tx != nil {
		res, err = tx.Exec(query, product.Name, product.Description, product.Category, product.Price, product.Id)
	} else {
		res, err = sc.db.Exec(query, product.Name, product.Description, product.Category, product.Price, product.Id)
	}
	if err != nil {
		return err
	}
	if a, err := res.RowsAffected(); err != nil {
		return err
	} else if a == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (sc *StoreContext) DeleteProduct(tx *sql.Tx, id int) error {
	query := "DELETE FROM product WHERE id = $1;"
	var res sql.Result
	var err error
	if tx != nil {
		res, err = tx.Exec(query, id)
	} else {
		res, err = sc.db.Exec(query, id)
	}
	if err != nil {
		return err
	}
	if a, err := res.RowsAffected(); err != nil {
		return err
	} else if a == 0 {
		return sql.ErrNoRows
	}
	return nil
}
