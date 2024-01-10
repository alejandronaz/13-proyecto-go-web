package main

import (
	"goweb/app/internals/handlers"
	"goweb/app/internals/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	// create the repo
	repo := repository.GetRepository()
	// load the data
	repo.LoadData()

	// create a router with chi
	router := chi.NewRouter()

	// create the routes
	router.Get("/ping", handlers.PingHandler)

	router.Route("/products", func(r chi.Router) {
		r.Get("/", handlers.GetAllProductsHandler)
		r.Get("/{id}", handlers.GetProductByIDHandler)
		r.Get("/search", handlers.GetProductsByPriceGreaterThanHandler)
		r.Post("/", handlers.CreateProductHandler)
	})

	http.ListenAndServe(":8080", router)

}
