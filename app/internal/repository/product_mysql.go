package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"goweb/app/internal"
	"os"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLConnection() (*sql.DB, error) {
	config := mysql.Config{
		User:      "root",
		Passwd:    os.Getenv("MYSQL_ROOT_PASSWORD"),
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "my_db",
		ParseTime: true,
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		fmt.Println("error connecting to the database: ", err)
		return nil, err
	}

	// check the connection
	if err = db.Ping(); err != nil {
		fmt.Println("error pinging the database: ", err)
		return nil, err
	}

	return db, nil
}

func NewProductRepositorySQL(db *sql.DB) *ProductRepositorySQL {
	return &ProductRepositorySQL{
		db: db,
	}
}

type ProductRepositorySQL struct {
	db *sql.DB
}

// GetAllProducts returns all products
func (r *ProductRepositorySQL) GetAllProducts() []internal.Product {

	rows, err := r.db.Query(
		"SELECT id, name, quantity, code_value, is_published, expiration, price FROM products",
	)

	if err != nil {
		fmt.Println("error querying the database: ", err)
		return nil
	}

	// iterate over the rows
	var products []internal.Product
	for rows.Next() {
		var product internal.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price)
		if err != nil {
			fmt.Println("error scanning the row: ", err)
			return nil
		}

		products = append(products, product)
	}

	return products
}

// GetProductByID returns a product by id
func (r *ProductRepositorySQL) GetProductByID(id int) internal.Product {

	// query
	row := r.db.QueryRow(
		"SELECT id, name, quantity, code_value, is_published, expiration, price FROM products WHERE id = ?",
		id,
	)

	if err := row.Err(); err != nil {
		fmt.Println("error querying the database: ", err)
		return internal.Product{}
	}

	var product internal.Product
	err := row.Scan(&product.ID, &product.Name, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Product not found")
			return product
		}
		return product
	}

	return product
}

// GetProductsByPriceGreaterThan returns products by price greater than
func (r *ProductRepositorySQL) GetProductsByPriceGreaterThan(price float64) []internal.Product {

	rows, err := r.db.Query(
		"SELECT id, name, quantity, code_value, is_published, expiration, price FROM products WHERE price > ?",
		price,
	)

	if err != nil {
		fmt.Println("error querying the database: ", err)
		return nil
	}

	// iterate over the rows
	var products []internal.Product
	for rows.Next() {
		var product internal.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price)
		if err != nil {
			fmt.Println("error scanning the row: ", err)
			return nil
		}

		products = append(products, product)
	}

	return products

}

// AddProduct adds a product
func (r *ProductRepositorySQL) AddProduct(product internal.Product) internal.Product {

	// query
	result, err := r.db.Exec(
		"INSERT INTO products (name, quantity, code_value, is_published, expiration, price) VALUES (?, ?, ?, ?, ?, ?)",
		product.Name, product.Quantity, product.CodeValue, product.IsPublished, product.Expiration, product.Price,
	)

	if err != nil {
		fmt.Println("error querying the database: ", err)
		return internal.Product{}
	}

	// get the id of the inserted product
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("error getting the last inserted id: ", err)
		return internal.Product{}
	}

	product.ID = int(id)
	return product

}

// UpdateProduct updates a product
func (r *ProductRepositorySQL) UpdateProduct(product internal.Product) (internal.Product, error) {

	// query
	_, err := r.db.Exec(
		"UPDATE products SET name = ?, quantity = ?, code_value = ?, is_published = ?, expiration = ?, price = ? WHERE id = ?",
		product.Name, product.Quantity, product.CodeValue, product.IsPublished, product.Expiration, product.Price, product.ID,
	)
	if err != nil {
		fmt.Println("error querying the database: ", err)
		return internal.Product{}, err
	}

	return product, nil
}

// DeleteProduct deletes a product
func (r *ProductRepositorySQL) DeleteProduct(id int) error {

	// query
	res, err := r.db.Exec(
		"DELETE FROM products WHERE id = ?",
		id,
	)
	if err != nil {
		fmt.Println("error querying the database: ", err)
		return err
	}

	// check if the product was deleted
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Println("error getting the rows affected: ", err)
		return err
	}
	if rowsAffected == 0 {
		fmt.Println("product not found")
		return errors.New("product not found")
	}

	return nil
}
