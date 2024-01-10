package handlers

import (
	"encoding/json"
	"errors"
	"goweb/app/internals/model"
	"goweb/app/internals/services"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {

	products := services.GetAllProducts()

	// // productsAsJSON is a slice of bytes
	// productsAsJSON, err := json.Marshal(products)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Header().Set("Content-Type", "text/plain")
	// 	w.Write([]byte("An error occurred"))
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if len(products) == 0 {
		w.Write([]byte(`[]`))
		return
	}

	// parse each product to ResponseBodyProduct
	productsAsResponse := parseProductsToBody(products)

	json.NewEncoder(w).Encode(productsAsResponse)
}

func GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the url
	id := chi.URLParam(r, "id")

	// convert the id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	product, err := services.GetProductByID(idInt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if errors.Is(err, services.ErrProductNotFound) {
		w.Write([]byte(`{}`))
		return
	}

	// parse product to ResponseBodyProduct
	productAsResponse := parseProductToBody(product)

	json.NewEncoder(w).Encode(productAsResponse)
}

func GetProductsByPriceGreaterThanHandler(w http.ResponseWriter, r *http.Request) {
	// get the price from the query param
	priceGt := r.URL.Query().Get("priceGt")

	// convert the price to float
	priceGtFloat, err := strconv.ParseFloat(priceGt, 64)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid price"))
		return
	}

	// get the products by price
	products := services.GetProductsByPriceGreaterThan(priceGtFloat)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if len(products) == 0 {
		w.Write([]byte(`[]`))
		return
	}

	// parse each product to ResponseBodyProduct
	productsAsResponse := parseProductsToBody(products)

	json.NewEncoder(w).Encode(productsAsResponse)

}

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {

	// -----------------------------------------------------
	// check auth header for creating a product
	// TODO: implement a middleware for this
	authHeader := r.Header.Get("Authorization")
	if authHeader != "1234" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}
	// -----------------------------------------------------

	// check if json sent by client has all the required fields:
	// 1. get the body as bytes
	bytesJson, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid product"))
		return
	}
	// 2. parse the bytes to a map (simil json)
	var bodyJson map[string]any
	err = json.Unmarshal(bytesJson, &bodyJson)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid product"))
		return
	}
	// 3. check if the map has all the required fields
	err = checkRequiredFields(bodyJson, "name", "quantity", "code_value", "is_published", "expiration", "price")
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// get the product from the request body
	var product RequestBodyProduct

	err = json.Unmarshal(bytesJson, &product)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid product"))
		return
	}

	// parse RequestBody to product model
	newProduct := model.Product{
		Name:        product.Name,
		Quantity:    product.Quantity,
		CodeValue:   product.CodeValue,
		IsPublished: product.IsPublished,
		Expiration:  product.Expiration,
		Price:       product.Price,
	}

	// create the product
	newProduct, err = services.CreateProduct(newProduct)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // no es recomendado retornar el error directamente, dado que puede exponer informacion interna
		return
	}

	// parse product to ResponseBodyProduct
	productAsResponse := parseProductToBody(newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(productAsResponse)

}
