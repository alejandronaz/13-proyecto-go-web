package application

import (
	"errors"
	"goweb/app/internal/handler"
	"goweb/app/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ServerChi struct {
	address string
}

func NewServer(address string) *ServerChi {
	defaultAddress := ":8080"
	if address == "" {
		address = defaultAddress
	}
	return &ServerChi{
		address: address,
	}
}

func (s *ServerChi) Start() error {

	// create the repo
	repo := repository.GetRepository()
	// load the data
	repo.LoadData()

	// create a router with chi
	router := chi.NewRouter()

	// create the routes
	router.Get("/ping", handler.PingHandler)

	router.Route("/products", func(r chi.Router) {
		r.Get("/", handler.GetAllProductsHandler)
		r.Get("/{id}", handler.GetProductByIDHandler)
		r.Get("/search", handler.GetProductsByPriceGreaterThanHandler)
		r.Post("/", handler.CreateProductHandler)
	})

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		return errors.New("An error occurred while starting the server")
	}
	return nil
}
