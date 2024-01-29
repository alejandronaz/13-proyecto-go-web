package application

import (
	"errors"
	"goweb/app/internal/handler"
	"goweb/app/internal/middleware"
	"goweb/app/internal/repository"
	"goweb/app/internal/service"
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

	// Initialize the dependencies
	// 1. Repository
	// 2. Service
	// 3. Handler

	// 0. create the db connection
	db, err := repository.NewMySQLConnection()
	if err != nil {
		return errors.New("an error occurred while connecting to the database")
	}
	defer db.Close()

	// 1. create the repo (aqui elijo especificamente que repo usar)
	repo := repository.NewProductRepositorySQL(db)
	// 2. create the service
	service := service.NewProductService(repo)
	// 3. create the handler
	handler := handler.NewProductHandler(service)

	// create a router with chi
	router := chi.NewRouter()

	// add the middleware
	router.Use(middleware.Logs, middleware.Auth)

	// create the routes
	router.Get("/ping", handler.Ping)

	router.Route("/products", func(r chi.Router) {
		r.Get("/", handler.GetAllProducts)
		r.Get("/{id}", handler.GetProductByID)
		r.Get("/search", handler.GetProductsByPriceGreaterThan)
		r.Post("/", handler.CreateProduct)
		r.Put("/{id}", handler.UpdateProduct)
		r.Patch("/{id}", handler.ParcialUpdateProduct)
		r.Delete("/{id}", handler.DeleteProduct)

		r.Get("/consumer_price", handler.CalculateConsumerPrice)
	})

	// 5. start the server
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		return errors.New("an error occurred while starting the server")
	}
	return nil
}
