package service

import (
	"goweb/app/internal"
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

	// add the product to the repo
	product = p.repo.AddProduct(product)

	return product, nil

}

func (p *ProductService) UpdateProduct(product internal.Product) (internal.Product, error) {

	// check if product is empty
	if product.IsEmpty() {
		return internal.Product{}, internal.ErrProductEmpty
	}

	// check if the code value belongs to another product
	products := p.repo.GetAllProducts()
	for _, p := range products {
		if p.CodeValue == product.CodeValue && p.ID != product.ID {
			return internal.Product{}, internal.ErrCodeValueBelongsToOther
		}
	}

	prodUpdt, err := p.repo.UpdateProduct(product)
	if err != nil {
		return internal.Product{}, err
	}

	return prodUpdt, nil

}

func (p *ProductService) DeleteProduct(id int) error {

	err := p.repo.DeleteProduct(id)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductService) CalculateConsumerPrice(idList ...int) ([]internal.Product, float64, error) {

	// calculate the number of each product in the list
	var idMap = make(map[int]int)
	for _, id := range idList {
		idMap[id]++
	}

	// if no id is passed, then calculate the price for all products
	if len(idList) == 0 {
		for _, prod := range p.GetAllProducts() {
			idMap[prod.ID]++
		}
	}

	// calculate the tax
	tax := 1.0
	switch {
	case len(idList) < 10:
		tax = 1.21
	case len(idList) > 10 && len(idList) <= 20:
		tax = 1.17
	case len(idList) > 20:
		tax = 1.15
	}

	finalPrice := 0.0
	prods := []internal.Product{}

	for id, quantity := range idMap {
		product, _ := p.GetProductByID(id)
		if product.Quantity >= quantity {
			finalPrice += product.Price * float64(quantity)
			product.Quantity = quantity // set the quantity requested by the consumer
			prods = append(prods, product)
		}
	}

	return prods, finalPrice * tax, nil

}
