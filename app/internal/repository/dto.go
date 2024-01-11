package repository

import "goweb/app/internal"

type ProductDTO struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func internalToDTO(products []internal.Product) []ProductDTO {

	var productsDTO []ProductDTO

	for _, product := range products {
		productsDTO = append(productsDTO, ProductDTO{
			ID:          product.ID,
			Name:        product.Name,
			CodeValue:   product.CodeValue,
			Expiration:  product.Expiration,
			IsPublished: product.IsPublished,
			Quantity:    product.Quantity,
			Price:       product.Price,
		})
	}

	return productsDTO
}

func dtoToInternal(products []ProductDTO) []internal.Product {

	var productsInternal []internal.Product

	for _, product := range products {
		productsInternal = append(productsInternal, internal.Product{
			ID:          product.ID,
			Name:        product.Name,
			CodeValue:   product.CodeValue,
			Expiration:  product.Expiration,
			IsPublished: product.IsPublished,
			Quantity:    product.Quantity,
			Price:       product.Price,
		})
	}

	return productsInternal

}
