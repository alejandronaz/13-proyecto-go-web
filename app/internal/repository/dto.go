package repository

import (
	"goweb/app/internal"
	"time"
)

type ProductDTO struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func internalsToDTOs(products []internal.Product) []ProductDTO {

	var productsDTO []ProductDTO

	for _, product := range products {
		productsDTO = append(productsDTO, ProductDTO{
			ID:          product.ID,
			Name:        product.Name,
			CodeValue:   product.CodeValue,
			Expiration:  product.Expiration.Format("02/01/2006"),
			IsPublished: product.IsPublished,
			Quantity:    product.Quantity,
			Price:       product.Price,
		})
	}

	return productsDTO
}

func dtosToInternals(products []ProductDTO) []internal.Product {

	var productsInternal []internal.Product

	for _, product := range products {
		parsedExpiration, _ := parseExpirationToTime(product.Expiration)
		productsInternal = append(productsInternal, internal.Product{
			ID:          product.ID,
			Name:        product.Name,
			CodeValue:   product.CodeValue,
			Expiration:  parsedExpiration,
			IsPublished: product.IsPublished,
			Quantity:    product.Quantity,
			Price:       product.Price,
		})
	}

	return productsInternal

}

func parseExpirationToTime(expiration string) (time.Time, error) {
	// if time cant parse it, then it is invalid
	parsedTime, err := time.Parse("02/01/2006", expiration)
	if err != nil {
		return time.Time{}, internal.ErrInvalidExpirationFormat
	}
	return parsedTime, nil
}
