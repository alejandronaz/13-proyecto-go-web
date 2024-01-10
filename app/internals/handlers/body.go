package handlers

import (
	"errors"
	"goweb/app/internals/model"
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

func parseProductToBody(product model.Product) ResponseBodyProduct {
	return ResponseBodyProduct{
		ID:          product.ID,
		Name:        product.Name,
		Quantity:    product.Quantity,
		CodeValue:   product.CodeValue,
		IsPublished: product.IsPublished,
		Expiration:  product.Expiration,
		Price:       product.Price,
	}
}

func parseProductsToBody(products []model.Product) []ResponseBodyProduct {
	var productsAsResponse []ResponseBodyProduct
	for _, product := range products {
		productsAsResponse = append(productsAsResponse, parseProductToBody(product))
	}
	return productsAsResponse
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
