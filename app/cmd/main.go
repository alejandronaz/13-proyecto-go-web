package main

import (
	"encoding/json"
	"goweb/app/internals/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// create the repo
var repo = repository.Repository{}

func main() {

	// load the data
	repo.LoadData()

	// create a router with chi
	router := chi.NewRouter()

	// create the routes
	router.Get("/ping", PingHandler)

	router.Route("/products", func(r chi.Router) {
		r.Get("/", GetAllProductsHandler)
		r.Get("/{id}", GetProductByIDHandler)
		r.Get("/search", GetProductsByPriceGreaterThanHandler)
	})

	http.ListenAndServe(":3000", router)

}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	products := repo.GetAllProducts()

	// productsAsJSON is a slice of bytes
	productsAsJSON, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(productsAsJSON)
}

func GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the url
	id := chi.URLParam(r, "id")

	// convert the id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	// get the product by id
	product := repo.GetProductByID(idInt)

	// productAsJSON is a slice of bytes
	productAsJSON, err := json.Marshal(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(productAsJSON)
}

func GetProductsByPriceGreaterThanHandler(w http.ResponseWriter, r *http.Request) {
	// get the price from the query param
	priceGt := r.URL.Query().Get("priceGt")

	// convert the price to float
	priceGtFloat, err := strconv.ParseFloat(priceGt, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid price"))
		return
	}

	// get the products by price
	products := repo.GetProductsByPriceGreaterThan(priceGtFloat)

	// parse to json
	productsAsJSON, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(productsAsJSON)

}
