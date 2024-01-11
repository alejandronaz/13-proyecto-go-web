package repository

import (
	"encoding/json"
	"fmt"
	"goweb/app/internal"
	"os"
)

// implements the ProductRepository interface
type RepositoryMap struct {
	Products map[int]internal.Product
	lastID   int
}

func NewRepositoryMap() *RepositoryMap {
	repo := &RepositoryMap{
		Products: make(map[int]internal.Product),
	}
	repo.LoadData()
	return repo
}

func (r *RepositoryMap) LoadData() {

	// read the json file as a slice of bytes
	data, err := os.ReadFile("app/data/products.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// unmarshal the bytes to a slice
	var products []ProductDTO
	err = json.Unmarshal([]byte(data), &products)
	if err != nil {
		fmt.Println("Hubo un error")
		return
	}

	// convert the slice of DTOs to a slice of internal.Product
	productsInternal := dtoToInternal(products)

	// convert the slice to a map
	lastId := 0
	for _, product := range productsInternal {
		r.Products[product.ID] = product
		if product.ID > lastId {
			lastId = product.ID
		}
	}
	r.lastID = lastId

}

// implement the methods from the interface internal.ProductRepository
func (r *RepositoryMap) GetAllProducts() []internal.Product {
	var products []internal.Product
	for _, product := range r.Products {
		products = append(products, product)
	}
	return products
}

func (r *RepositoryMap) GetProductByID(id int) internal.Product {

	prod, ok := r.Products[id]
	if !ok {
		return internal.Product{}
	}

	return prod

}

func (r *RepositoryMap) GetProductsByPriceGreaterThan(price float64) []internal.Product {
	var products []internal.Product
	for _, product := range r.Products {
		if product.Price > price {
			products = append(products, product)
		}
	}
	return products
}

func (r *RepositoryMap) AddProduct(product internal.Product) internal.Product {
	r.lastID++
	product.ID = r.lastID
	r.Products[r.lastID] = product

	return product
}

func (r *RepositoryMap) UpdateProduct(product internal.Product) (internal.Product, error) {

	for id, prod := range r.Products {

		if id == product.ID {
			prod.Name = product.Name
			prod.CodeValue = product.CodeValue
			prod.Expiration = product.Expiration
			prod.IsPublished = product.IsPublished
			prod.Quantity = product.Quantity
			prod.Price = product.Price

			r.Products[id] = prod

			return prod, nil
		}
	}

	return internal.Product{}, internal.ErrProductNotFound
}

func (r *RepositoryMap) DeleteProduct(id int) error {
	delete(r.Products, id)
	return nil
}
