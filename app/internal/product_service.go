package internal

import "errors"

type ProductService interface {
	GetAllProducts() []Product
	GetProductByID(id int) (Product, error)
	GetProductsByPriceGreaterThan(price float64) []Product
	CreateProduct(product Product) (Product, error)
	UpdateOrCreateProduct(product Product) (Product, error)
	// for patch
	UpdateProduct(product Product) (Product, error)
	DeleteProduct(id int) error
}

var (
	ErrProductNotFound         = errors.New("product not found")
	ErrProductExists           = errors.New("product already exists")
	ErrProductEmpty            = errors.New("product is empty")
	ErrInvalidExpirationFormat = errors.New("invalid expiration format")
	ErrCodeValueBelongsToOther = errors.New("code value belongs to other product")
)
