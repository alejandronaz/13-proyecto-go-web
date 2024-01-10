package handler

import (
	"encoding/json"
	"errors"
	"goweb/app/internal"
	"goweb/app/internal/platform/web/response"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	service internal.ProductService
}

func NewProductHandler(service internal.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (p *ProductHandler) Ping() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		response.Text(w, http.StatusOK, "pong")
	}
}

func (p *ProductHandler) GetAllProducts() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		products := p.service.GetAllProducts()

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
}

func (p *ProductHandler) GetProductByID() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the id from the url
		id := chi.URLParam(r, "id")

		// convert the id to int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		product, err := p.service.GetProductByID(idInt)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if errors.Is(err, internal.ErrProductNotFound) {
			w.Write([]byte(`{}`))
			return
		}

		// parse product to ResponseBodyProduct
		productAsResponse := parseProductToBody(product)

		json.NewEncoder(w).Encode(productAsResponse)
	}
}

func (p *ProductHandler) GetProductsByPriceGreaterThan() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// get the price from the query param
		priceGt := r.URL.Query().Get("priceGt")

		// convert the price to float
		priceGtFloat, err := strconv.ParseFloat(priceGt, 64)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "Invalid price")
			return
		}

		// get the products by price
		products := p.service.GetProductsByPriceGreaterThan(priceGtFloat)

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
}

func (p *ProductHandler) CreateProduct() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// -----------------------------------------------------
		// check auth header for creating a product
		// TODO: implement a middleware for this
		authHeader := r.Header.Get("Authorization")
		if authHeader != "1234" {
			response.Text(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		// -----------------------------------------------------

		// --------check if json sent by client has all the required fields--------
		// 1. get the body as bytes
		bytesJson, err := io.ReadAll(r.Body)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "Invalid product")
			return
		}
		// 2. parse the bytes to a map (simil json)
		var bodyJson map[string]any
		err = json.Unmarshal(bytesJson, &bodyJson)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "Invalid product")
			return
		}
		// 3. check if the map has all the required fields
		err = checkRequiredFields(bodyJson, "name", "quantity", "code_value", "is_published", "expiration", "price")
		if err != nil {
			response.Text(w, http.StatusBadRequest, err.Error())
			return
		}
		// -------------------------------------------------------------------------

		// get the product from the request body
		var product RequestBodyProduct

		err = json.Unmarshal(bytesJson, &product)
		if err != nil {
			response.Text(w, http.StatusBadRequest, "Invalid product")
			return
		}

		// parse RequestBody to product model
		newProduct := internal.Product{
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}

		// create the product
		newProduct, err = p.service.CreateProduct(newProduct)
		if err != nil {
			// no es recomendado retornar el error directamente, dado que puede exponer informacion interna
			response.Text(w, http.StatusBadRequest, err.Error())
			return
		}

		// parse product to ResponseBodyProduct
		productAsResponse := parseProductToBody(newProduct)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(productAsResponse)

	}
}
