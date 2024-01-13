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

func NewRepository(data []internal.Product) *Repository {

	if data == nil {
		repo := &Repository{}
		repo.LoadData()
		return repo
	}

	return &Repository{
		Products: data,
	}
}

func (r *Repository) LoadData() {

	// read the json file as a slice of bytes
	data, err := os.ReadFile("app/data/products.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// unmarshal the bytes to the repo slice
	var products []ProductDTO
	err = json.Unmarshal([]byte(data), &products)
	if err != nil {
		fmt.Println("Hubo un error")
		return
	}

	// convert the slice of DTOs to a slice of internal.Product
	r.Products = dtosToInternals(products)

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

func (r *Repository) AddProduct(product internal.Product) internal.Product {

	product.ID = len(r.Products) + 1
	r.Products = append(r.Products, product)

	return product
}

func (r *Repository) UpdateProduct(product internal.Product) (internal.Product, error) {
	for i, p := range r.Products {
		if p.ID == product.ID {
			r.Products[i].Name = product.Name
			r.Products[i].CodeValue = product.CodeValue
			r.Products[i].Expiration = product.Expiration
			r.Products[i].IsPublished = product.IsPublished
			r.Products[i].Quantity = product.Quantity
			r.Products[i].Price = product.Price
			return r.Products[i], nil
		}
	}

	return internal.Product{}, internal.ErrProductNotFound
}

func (r *Repository) DeleteProduct(id int) error {
	for i, p := range r.Products {
		if p.ID == id {
			r.Products = append(r.Products[:i], r.Products[i+1:]...)
			return nil
		}
	}

	return internal.ErrProductNotFound
}
