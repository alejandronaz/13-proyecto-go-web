package repository

import (
	"encoding/json"
	"fmt"
	"goweb/app/internals/model"
	"os"
)

type Repository struct {
	Products []model.Product
}

var repo = Repository{}

func GetRepository() *Repository {
	return &repo
}

func (r *Repository) LoadData() {

	// read the json file as a slice of bytes
	data, err := os.ReadFile("app/data/products.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// unmarshal the bytes to the repo slice
	err = json.Unmarshal([]byte(data), &r.Products)
	if err != nil {
		fmt.Println("Hubo un error")
		return
	}

}

func (r *Repository) GetAllProducts() []model.Product {
	return r.Products
}

func (r *Repository) GetProductByID(id int) model.Product {
	for _, product := range r.Products {
		if product.ID == id {
			return product
		}
	}
	return model.Product{}
}

func (r *Repository) GetProductsByPriceGreaterThan(price float64) []model.Product {
	var products []model.Product
	for _, product := range r.Products {
		if product.Price > price {
			products = append(products, product)
		}
	}
	return products
}
