package handlers

import (
	"encoding/json"
	"goweb/app/internals/repository"
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

	// get the repo
	repo := repository.GetRepository()

	products := repo.GetAllProducts()

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
	json.NewEncoder(w).Encode(products)
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

	// get the repo
	repo := repository.GetRepository()

	// get the product by id
	product := repo.GetProductByID(idInt)

	// // productAsJSON is a slice of bytes
	// productAsJSON, err := json.Marshal(product)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Header().Set("Content-Type", "text/plain")
	// 	w.Write([]byte("An error occurred"))
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if product.IsEmpty() {
		w.Write([]byte(`{}`))
		return
	}

	json.NewEncoder(w).Encode(product)
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

	// get the repo
	repo := repository.GetRepository()

	// get the products by price
	products := repo.GetProductsByPriceGreaterThan(priceGtFloat)

	// // parse to json
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
	json.NewEncoder(w).Encode(products)

}
