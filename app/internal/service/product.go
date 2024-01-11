package service

import (
	"fmt"
	"goweb/app/internal"
	"strings"
	"time"
)

// implements internal.ProductService and uses internal.ProductRepository (other interface)
type ProductService struct {
	repo internal.ProductRepository
}

// create a new product service, which uses a product repository passed through the constructor
func NewProductService(repo internal.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// implement the methods from the interface internal.ProductService
func (p *ProductService) GetAllProducts() []internal.Product {
	return p.repo.GetAllProducts()
}

func (p *ProductService) GetProductByID(id int) (internal.Product, error) {

	product := p.repo.GetProductByID(id)

	if product.IsEmpty() {
		return product, internal.ErrProductNotFound
	}

	return product, nil
}

func (p *ProductService) GetProductsByPriceGreaterThan(price float64) []internal.Product {
	return p.repo.GetProductsByPriceGreaterThan(price)
}

func (p *ProductService) CreateProduct(product internal.Product) (internal.Product, error) {

	products := p.repo.GetAllProducts()

	// check if product is empty
	if product.IsEmpty() {
		return internal.Product{}, internal.ErrProductEmpty
	}

	// check if the value_code already exists
	for _, p := range products {
		if p.CodeValue == product.CodeValue {
			return internal.Product{}, internal.ErrProductExists
		}
	}

	// verify expiration format XX/XX/XXXX
	exp := strings.Split(product.Expiration, "/")
	if len(exp) != 3 {
		return internal.Product{}, internal.ErrInvalidExpirationFormat
	}
	// if time cant parse it, then it is invalid
	_, err := time.Parse(time.DateOnly, fmt.Sprint(exp[2], "-", exp[1], "-", exp[0]))
	if err != nil {
		return internal.Product{}, internal.ErrInvalidExpirationFormat
	}

	// set id to the product
	product.ID = len(products) + 1

	// add the product to the repo
	p.repo.AddProduct(product)

	return product, nil

}

func (p *ProductService) UpdateOrCreateProduct(product internal.Product) (internal.Product, error) {

	// 1. Update

	// check if product is empty
	if product.IsEmpty() {
		return internal.Product{}, internal.ErrProductEmpty
	}

	// verify expiration format XX/XX/XXXX
	exp := strings.Split(product.Expiration, "/")
	if len(exp) != 3 {
		return internal.Product{}, internal.ErrInvalidExpirationFormat
	}
	// if time cant parse it, then it is invalid
	_, err := time.Parse(time.DateOnly, fmt.Sprint(exp[2], "-", exp[1], "-", exp[0]))
	if err != nil {
		return internal.Product{}, internal.ErrInvalidExpirationFormat
	}

	// check if the code value belongs to another product
	products := p.repo.GetAllProducts()
	for _, p := range products {
		if p.CodeValue == product.CodeValue && p.ID != product.ID {
			return internal.Product{}, internal.ErrCodeValueBelongsToOther
		}
	}

	prodUpdt, err := p.repo.UpdateProduct(product)
	if err == nil { // means that the product was updated
		return prodUpdt, nil
	}

	// 2. Create
	newProd := internal.Product{
		ID:          len(products) + 1,
		Name:        product.Name,
		CodeValue:   product.CodeValue,
		Price:       product.Price,
		Expiration:  product.Expiration,
		Quantity:    product.Quantity,
		IsPublished: product.IsPublished,
	}
	p.repo.AddProduct(newProd)
	return newProd, nil

}
