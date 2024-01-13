package repository

import (
	"encoding/json"
	"fmt"
	"goweb/app/internal"
	"os"
)

// implements the ProductRepository interface
type RepositoryFile struct {
	lastID int
}

func NewRepositoryFile() *RepositoryFile {
	return &RepositoryFile{}
}

func (r *RepositoryFile) getDataFromFile() ([]internal.Product, error) {

	// read the json file as a slice of bytes
	data, err := os.ReadFile("app/data/file_storage/products.json")
	if err != nil {
		fmt.Println(err)
		return []internal.Product{}, err
	}

	// unmarshal the bytes to a slice
	var products []ProductDTO
	err = json.Unmarshal([]byte(data), &products)
	if err != nil {
		return []internal.Product{}, err
	}

	// get the last id
	lastId := 0
	for _, product := range products {
		if product.ID > lastId {
			lastId = product.ID
		}
	}
	r.lastID = lastId

	// convert the slice of DTOs to a slice of internal.Product
	productsInternal := dtosToInternals(products)

	return productsInternal, nil

}

func (r *RepositoryFile) saveDataToFile(products []internal.Product) error {

	// open the file
	file, err := os.OpenFile("app/data/file_storage/products.json", os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	// parse internal.Product[] to ProductDTO
	productsDTO := internalsToDTOs(products)

	// convert the slice to bytes
	bytes, err := json.Marshal(productsDTO)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// write the bytes to the file
	_, err = file.Write(bytes)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

// implement the methods from the interface internal.ProductRepository
func (r *RepositoryFile) GetAllProducts() []internal.Product {

	products, _ := r.getDataFromFile()
	return products
}

func (r *RepositoryFile) GetProductByID(id int) internal.Product {

	products, _ := r.getDataFromFile()

	for _, product := range products {
		if product.ID == id {
			return product
		}
	}

	return internal.Product{}

}

func (r *RepositoryFile) GetProductsByPriceGreaterThan(price float64) []internal.Product {

	products, _ := r.getDataFromFile()

	var productsSorted []internal.Product

	for _, product := range products {
		if product.Price > price {
			productsSorted = append(productsSorted, product)
		}
	}
	return productsSorted
}

func (r *RepositoryFile) AddProduct(product internal.Product) internal.Product {

	products, _ := r.getDataFromFile()

	r.lastID++
	product.ID = r.lastID

	products = append(products, product)

	err := r.saveDataToFile(products)
	if err != nil {
		fmt.Println(err)
	}

	return product

}

func (r *RepositoryFile) UpdateProduct(product internal.Product) (internal.Product, error) {

	products, _ := r.getDataFromFile()

	for i, prod := range products {

		if prod.ID == product.ID {
			prod.Name = product.Name
			prod.CodeValue = product.CodeValue
			prod.Expiration = product.Expiration
			prod.IsPublished = product.IsPublished
			prod.Quantity = product.Quantity
			prod.Price = product.Price

			products[i] = prod

			r.saveDataToFile(products)

			return prod, nil
		}
	}

	return internal.Product{}, internal.ErrProductNotFound
}

func (r *RepositoryFile) DeleteProduct(id int) error {

	products, _ := r.getDataFromFile()

	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)

			r.saveDataToFile(products)

			return nil
		}
	}

	return internal.ErrProductNotFound
}
