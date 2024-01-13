package handler

import (
	"errors"
	"goweb/app/internal"
	"time"
)

type RequestBodyProduct struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type ResponseBodyProduct struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func parseProductToBody(product internal.Product) ResponseBodyProduct {

	return ResponseBodyProduct{
		ID:          product.ID,
		Name:        product.Name,
		Quantity:    product.Quantity,
		CodeValue:   product.CodeValue,
		IsPublished: product.IsPublished,
		Expiration:  product.Expiration.Format("02/01/2006"),
		Price:       product.Price,
	}
}

func parseProductsToBody(products []internal.Product) []ResponseBodyProduct {
	var productsAsResponse []ResponseBodyProduct
	for _, product := range products {
		productsAsResponse = append(productsAsResponse, parseProductToBody(product))
	}
	return productsAsResponse
}

func parseBodyToProduct(id int, body RequestBodyProduct) (internal.Product, error) {
	// if time cant parse it, then it is invalid
	parsedTime, err := time.Parse("02/01/2006", body.Expiration)
	if err != nil {
		return internal.Product{}, internal.ErrInvalidExpirationFormat
	}
	return internal.Product{
		ID:          id,
		Name:        body.Name,
		Quantity:    body.Quantity,
		CodeValue:   body.CodeValue,
		IsPublished: body.IsPublished,
		Expiration:  parsedTime,
		Price:       body.Price,
	}, nil
}

func checkRequiredFields(body map[string]any, requiredFields ...string) error {
	for _, field := range requiredFields {
		_, ok := body[field]
		if !ok {
			return errors.New("Missing field " + field)
		}
	}
	return nil
}

type ResponseConsumerPrice struct {
	Products   []ResponseBodyProduct `json:"products"`
	TotalPrice float64               `json:"total_price"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
