package handlers

import (
	"encoding/json"
	"errors"
	"goweb/app/internals/model"
	"goweb/app/internals/services"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

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
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Invalid ID"))
		return
	}

	product, err := services.GetProductByID(idInt)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

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
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Invalid price"))
		return
	}

	// get the products by price
	products := services.GetProductsByPriceGreaterThan(priceGtFloat)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if len(products) == 0 {
		w.Write([]byte(`[]`))
		return
	}

	// parse each product to ResponseBodyProduct
	productsAsResponse := parseProductsToBody(products)

	json.NewEncoder(w).Encode(productsAsResponse)

}

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {

	// get the product from the request body
	var product RequestBodyProduct

	// IMPORTANTE: si alguna clave no coincide con el struct, se creará el struct con campos con zero values
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
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
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(err.Error()))
		return
	}

	// parse product to ResponseBodyProduct
	productAsResponse := parseProductToBody(newProduct)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productAsResponse)

}
