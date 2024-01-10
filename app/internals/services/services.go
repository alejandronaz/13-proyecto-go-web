package services

import (
	"fmt"
	"goweb/app/internals/model"
	"goweb/app/internals/repository"
	"strings"
	"time"
)

func GetAllProducts() []model.Product {
	// get the repo
	repo := repository.GetRepository()
	return repo.GetAllProducts()
}

func GetProductByID(id int) (model.Product, error) {
	// get the repo
	repo := repository.GetRepository()
	product := repo.GetProductByID(id)

	if product.IsEmpty() {
		return product, ErrProductNotFound
	}

	return product, nil
}

func GetProductsByPriceGreaterThan(price float64) []model.Product {
	// get the repo
	repo := repository.GetRepository()
	return repo.GetProductsByPriceGreaterThan(price)
}

func CreateProduct(product model.Product) (model.Product, error) {

	repo := repository.GetRepository()
	products := repo.GetAllProducts()

	// check if product is empty
	if product.IsEmpty() {
		return model.Product{}, ErrProductEmpty
	}

	// check if the value_code already exists
	for _, p := range products {
		if p.CodeValue == product.CodeValue {
			return model.Product{}, ErrProductExists
		}
	}

	// verify expiration format XX/XX/XXXX
	exp := strings.Split(product.Expiration, "/")
	if len(exp) != 3 {
		return model.Product{}, ErrInvalidExpirationFormat
	}
	// if time cant parse it, then it is invalid
	_, err := time.Parse(time.DateOnly, fmt.Sprint(exp[2], "-", exp[1], "-", exp[0]))
	if err != nil {
		return model.Product{}, ErrInvalidExpirationFormat
	}

	// set id to the product
	product.ID = len(products) + 1

	// add the product to the repo
	repo.AddProduct(product)

	return product, nil

}
