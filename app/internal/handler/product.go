package handler

import (
	"encoding/json"
	"errors"
	"goweb/app/internal"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/bootcamp-go/web/response"
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

func (p *ProductHandler) Ping(w http.ResponseWriter, r *http.Request) {

	response.Text(w, http.StatusOK, "pong")

}

func (p *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {

	products := p.service.GetAllProducts()

	if len(products) == 0 {
		response.JSON(w, http.StatusNotFound, ErrorResponse{
			Message: "No products found",
			Status:  http.StatusNotFound,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// parse each product to ResponseBodyProduct
	productsAsResponse := parseProductsToBody(products)

	json.NewEncoder(w).Encode(productsAsResponse)

}

func (p *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {

	// get the id from the url
	id := chi.URLParam(r, "id")

	// convert the id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid ID",
			Status:  http.StatusBadRequest,
		})
		return
	}

	product, err := p.service.GetProductByID(idInt)

	if errors.Is(err, internal.ErrProductNotFound) {
		response.JSON(w, http.StatusNotFound, ErrorResponse{
			Message: "No products found",
			Status:  http.StatusNotFound,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// parse product to ResponseBodyProduct
	productAsResponse := parseProductToBody(product)

	json.NewEncoder(w).Encode(productAsResponse)
}

func (p *ProductHandler) GetProductsByPriceGreaterThan(w http.ResponseWriter, r *http.Request) {

	// get the price from the query param
	priceGt := r.URL.Query().Get("priceGt")

	// convert the price to float
	priceGtFloat, err := strconv.ParseFloat(priceGt, 64)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid price",
			Status:  http.StatusBadRequest,
		})
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

func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {

	// --------check if json sent by client has all the required fields--------
	// 1. get the body as bytes
	bytesJson, err := io.ReadAll(r.Body)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid product",
			Status:  http.StatusBadRequest,
		})
		return
	}
	// 2. parse the bytes to a map (simil json)
	var bodyJson map[string]any
	err = json.Unmarshal(bytesJson, &bodyJson)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid product",
			Status:  http.StatusBadRequest,
		})
		return
	}
	// 3. check if the map has all the required fields
	err = checkRequiredFields(bodyJson, "name", "quantity", "code_value", "is_published", "expiration", "price")
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "There are missing fields",
			Status:  http.StatusBadRequest,
		})
		return
	}
	// -------------------------------------------------------------------------

	// get the product from the request body
	var product RequestBodyProduct

	err = json.Unmarshal(bytesJson, &product)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid product",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// parse RequestBody to product model
	newProduct, err := parseBodyToProduct(0, product)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrInvalidExpirationFormat):
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Invalid expiration format",
				Status:  http.StatusBadRequest,
			})
		default:
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Invalid product",
				Status:  http.StatusBadRequest,
			})
		}
		return
	}

	// create the product
	newProduct, err = p.service.CreateProduct(newProduct)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrProductExists):
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Product already exists",
				Status:  http.StatusBadRequest,
			})
		case errors.Is(err, internal.ErrInvalidExpirationFormat):
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Invalid expiration format",
				Status:  http.StatusBadRequest,
			})
		default:
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Invalid product",
				Status:  http.StatusBadRequest,
			})
		}
		return
	}

	// parse product to ResponseBodyProduct
	productAsResponse := parseProductToBody(newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(productAsResponse)

}

// for PUT method
func (p *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	// convert the id to int
	idProd, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid ID",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// --------check if json sent by client has all the required fields--------
	bytesJson, err := io.ReadAll(r.Body)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid product",
			Status:  http.StatusBadRequest,
		})
		return
	}
	var mapJson map[string]any
	if err := json.Unmarshal(bytesJson, &mapJson); err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid product",
			Status:  http.StatusBadRequest,
		})
		return
	}
	err = checkRequiredFields(mapJson, "name", "quantity", "code_value", "is_published", "expiration", "price")
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "There are missing fields",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// get the product from the request body
	var product RequestBodyProduct
	if err := json.Unmarshal(bytesJson, &product); err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid product",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// parse RequestBody to product model
	productModel, err := parseBodyToProduct(idProd, product)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrInvalidExpirationFormat):
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Invalid expiration format",
				Status:  http.StatusBadRequest,
			})
		default:
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Invalid product",
				Status:  http.StatusBadRequest,
			})
		}
		return
	}

	// call service
	productModel, err = p.service.UpdateProduct(productModel)
	if err != nil {
		response.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	// parse productModel to ResponseBodyProduct
	prodRes := parseProductToBody(productModel)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prodRes)

}

// for PATCH method
func (p *ProductHandler) ParcialUpdateProduct(w http.ResponseWriter, r *http.Request) {

	// convert the id to int
	idProd, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid ID",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// get the product by id
	product, err := p.service.GetProductByID(idProd)
	if errors.Is(err, internal.ErrProductNotFound) {
		response.JSON(w, http.StatusNotFound, ErrorResponse{
			Message: "No products found",
			Status:  http.StatusNotFound,
		})
		return
	}

	// get the product from the request body, setting the default values from the existent product
	var productBody RequestBodyProduct = RequestBodyProduct{
		Name:        product.Name,
		Quantity:    product.Quantity,
		CodeValue:   product.CodeValue,
		IsPublished: product.IsPublished,
		Expiration:  product.Expiration.Format("02/01/2006"),
		Price:       product.Price,
	}
	if err := json.NewDecoder(r.Body).Decode(&productBody); err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid product",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// parse RequestBody to product model
	productModel, err := parseBodyToProduct(idProd, productBody)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrInvalidExpirationFormat):
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Invalid expiration format",
				Status:  http.StatusBadRequest,
			})
		default:
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Invalid product",
				Status:  http.StatusBadRequest,
			})
		}
		return
	}

	// call service
	productModel, err = p.service.UpdateProduct(productModel)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrCodeValueBelongsToOther):
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Code value belongs to other product",
				Status:  http.StatusBadRequest,
			})
		case errors.Is(err, internal.ErrInvalidExpirationFormat):
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Invalid expiration format",
				Status:  http.StatusBadRequest,
			})
		case errors.Is(err, internal.ErrProductEmpty):
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Product is empty",
				Status:  http.StatusBadRequest,
			})
		default:
			response.JSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Invalid product",
				Status:  http.StatusBadRequest,
			})
		}
		return
	}

	// parse productModel to ResponseBodyProduct
	prodRes := parseProductToBody(productModel)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prodRes)

}

func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	// -----------------------------------------------------
	// check auth header for deleting a product
	authHeader := r.Header.Get("Authorization")
	if authHeader != "1234" {
		response.JSON(w, http.StatusUnauthorized, ErrorResponse{
			Message: "Unauthorized",
			Status:  http.StatusUnauthorized,
		})
		return
	}
	// -----------------------------------------------------

	// convert the id to int
	idProd, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Invalid ID",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// delete prod
	err = p.service.DeleteProduct(idProd)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "There was a problem deleting the product",
			Status:  http.StatusBadRequest,
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (p *ProductHandler) CalculateConsumerPrice(w http.ResponseWriter, r *http.Request) {

	list := r.URL.Query().Get("list")
	var sliceInt []int

	if list != "" {
		// convert to slice of int
		list = strings.ReplaceAll(list, "[", "")
		list = strings.ReplaceAll(list, "]", "")
		sliceStr := strings.Split(list, ",")

		for _, str := range sliceStr {
			integer, err := strconv.Atoi(str)
			if err != nil {
				response.JSON(w, http.StatusBadRequest, ErrorResponse{
					Message: "Invalid list",
					Status:  http.StatusBadRequest,
				})
				return
			}
			sliceInt = append(sliceInt, integer)
		}
	}

	// call service
	products, price, err := p.service.CalculateConsumerPrice(sliceInt...) // if no params are passed, sliceInt is empty, i.e. CalculateConsumerPrice()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "There was a problem calculating the consumer price",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// parse products to ResponseBodyProduct
	productsAsResponse := parseProductsToBody(products)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseConsumerPrice{
		Products:   productsAsResponse,
		TotalPrice: price,
	})

}
