package handler

import (
	"errors"
	"fmt"
	"goweb/app/internal"
	"strings"
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

func parseExpirationToTime(expiration string) (time.Time, error) {
	// verify expiration format XX/XX/XXXX
	exp := strings.Split(expiration, "/")
	if len(exp) != 3 {
		return time.Time{}, internal.ErrInvalidExpirationFormat
	}
	// if time cant parse it, then it is invalid
	parsedTime, err := time.Parse(time.DateOnly, fmt.Sprint(exp[2], "-", exp[1], "-", exp[0]))
	if err != nil {
		return time.Time{}, internal.ErrInvalidExpirationFormat
	}
	return parsedTime, nil
}

func parseProductsToBody(products []internal.Product) []ResponseBodyProduct {
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

type ResponseConsumerPrice struct {
	Products   []ResponseBodyProduct `json:"products"`
	TotalPrice float64               `json:"total_price"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
