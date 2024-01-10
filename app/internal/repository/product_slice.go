package repository

import (
	"encoding/json"
	"fmt"
	"goweb/app/internal"
	"os"
)

// implements the ProductRepository interface
type Repository struct {
	Products []internal.Product
}

func NewRepository() *Repository {
	repo := &Repository{}
	repo.LoadData()
	return repo
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

// implement the methods from the interface internal.ProductRepository
func (r *Repository) GetAllProducts() []internal.Product {
	return r.Products
}

func (r *Repository) GetProductByID(id int) internal.Product {
	for _, product := range r.Products {
		if product.ID == id {
			return product
		}
	}
	return internal.Product{}
}

func (r *Repository) GetProductsByPriceGreaterThan(price float64) []internal.Product {
	var products []internal.Product
	for _, product := range r.Products {
		if product.Price > price {
			products = append(products, product)
		}
	}
	return products
}

func (r *Repository) AddProduct(product internal.Product) {
	r.Products = append(r.Products, product)
}
