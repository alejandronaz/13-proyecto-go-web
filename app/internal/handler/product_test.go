package handler_test

import (
	"context"
	"goweb/app/internal"
	"goweb/app/internal/handler"
	"goweb/app/internal/repository"
	"goweb/app/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestGetAllProducts(t *testing.T) {

	t.Run("Se espera obtener todos los productos guardados.", func(t *testing.T) {
		// Arrange
		data := map[int]internal.Product{
			1: {
				ID:          1,
				Name:        "Producto 1",
				Quantity:    10,
				CodeValue:   "123456",
				IsPublished: true,
				Expiration:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
				Price:       100,
			},
			2: {
				ID:          2,
				Name:        "Producto 2",
				Quantity:    20,
				CodeValue:   "654321",
				IsPublished: false,
				Expiration:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
				Price:       200,
			},
		}
		repo := repository.NewRepositoryMap(data)
		service := service.NewProductService(repo)
		handler := handler.NewProductHandler(service)

		req := httptest.NewRequest("GET", "/products", nil)
		res := httptest.NewRecorder()

		// Act
		handler.GetAllProducts(res, req)

		// Assert
		expectedCode := http.StatusOK
		expectedBody := `[
							{"id":1,"name":"Producto 1","quantity":10,"code_value":"123456","is_published":true,"expiration":"31/12/2021","price":100},
							{"id":2,"name":"Producto 2","quantity":20,"code_value":"654321","is_published":false,"expiration":"31/12/2021","price":200}
						]`
		expectedHeader := http.Header{
			"Content-Type": []string{"application/json"},
		}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})

}

func TestGetProductByID(t *testing.T) {
	// test using path param
	t.Run("Obtener el producto con el id solicitado.", func(t *testing.T) {
		// Arrange
		data := map[int]internal.Product{
			1: {
				ID:          1,
				Name:        "Producto 1",
				Quantity:    10,
				CodeValue:   "123456",
				IsPublished: true,
				Expiration:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
				Price:       100,
			},
			2: {
				ID:          2,
				Name:        "Producto 2",
				Quantity:    20,
				CodeValue:   "654321",
				IsPublished: false,
				Expiration:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
				Price:       200,
			},
		}
		repo := repository.NewRepositoryMap(data)
		service := service.NewProductService(repo)
		handler := handler.NewProductHandler(service)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/products/1", nil)
		// add id path param to the request
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		context := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
		req = req.WithContext(context)

		// Act
		handler.GetProductByID(res, req)

		// Assert
		expectedCode := http.StatusOK
		expectedBody := `{"id":1,"name":"Producto 1","quantity":10,"code_value":"123456","is_published":true,"expiration":"31/12/2021","price":100}`
		expectedHeader := http.Header{
			"Content-Type": []string{"application/json"},
		}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})

	t.Run("El cliente envía un id erróneo", func(t *testing.T) {
		// Arrange
		data := map[int]internal.Product{}
		repo := repository.NewRepositoryMap(data)
		service := service.NewProductService(repo)
		handler := handler.NewProductHandler(service)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/products/NOTNUMBER", nil)
		// add id path param to the request
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "NOTNUMBER")
		context := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
		req = req.WithContext(context)

		// Act
		handler.GetProductByID(res, req)

		// Assert
		expectedCode := http.StatusBadRequest
		expectedBody := `{"message":"Invalid ID", "status": 400}`
		expectedHeader := http.Header{
			"Content-Type": []string{"application/json"},
		}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())

	})

	t.Run("Se solicita un id que no existe en los productos listados.", func(t *testing.T) {
		// Arrange
		data := map[int]internal.Product{
			1: {
				ID:          1,
				Name:        "Producto 1",
				Quantity:    10,
				CodeValue:   "123456",
				IsPublished: true,
				Expiration:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
				Price:       100,
			},
			2: {
				ID:          2,
				Name:        "Producto 2",
				Quantity:    20,
				CodeValue:   "654321",
				IsPublished: false,
				Expiration:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
				Price:       200,
			},
		}
		repo := repository.NewRepositoryMap(data)
		service := service.NewProductService(repo)
		handler := handler.NewProductHandler(service)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/products/3", nil)
		// add id path param to the request
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "3")
		context := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
		req = req.WithContext(context)

		// Act
		handler.GetProductByID(res, req)

		// Assert
		expectedCode := http.StatusNotFound
		expectedBody := `{"message":"No products found", "status": 404}`
		expectedHeader := http.Header{
			"Content-Type": []string{"application/json"},
		}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})
}

func TestCreateProduct(t *testing.T) {
	t.Run("Se añade un producto en la API y se devuelve el mismo en el cuerpo de la respuesta.", func(t *testing.T) {
		// Arrange
		data := map[int]internal.Product{}
		repo := repository.NewRepositoryMap(data)
		service := service.NewProductService(repo)
		handler := handler.NewProductHandler(service)

		body := strings.NewReader(`{"name":"Producto 1","quantity":10,"code_value":"123456","is_published":true,"expiration":"31/12/2021","price":100}`)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/products", body)
		req.Header.Add("Authorization", "1234")

		// Act
		handler.CreateProduct(res, req)

		// Assert
		expectedCode := http.StatusCreated
		expectedBody := `{"id":1,"name":"Producto 1","quantity":10,"code_value":"123456","is_published":true,"expiration":"31/12/2021","price":100}`
		expectedHeader := http.Header{
			"Content-Type": []string{"application/json"},
		}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("Se elimina el producto con dicho id, y no es necesario retornar nada.", func(t *testing.T) {
		// Arrange
		data := map[int]internal.Product{
			1: {
				ID:          1,
				Name:        "Producto 1",
				Quantity:    10,
				CodeValue:   "123456",
				IsPublished: true,
				Expiration:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
				Price:       100,
			},
		}
		repo := repository.NewRepositoryMap(data)
		service := service.NewProductService(repo)
		handler := handler.NewProductHandler(service)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/products/1", nil)
		req.Header.Add("Authorization", "1234")

		// add id path param to the request
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		context := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
		req = req.WithContext(context)

		// Act
		handler.DeleteProduct(res, req)

		// Assert
		expectedCode := http.StatusNoContent
		expectedBody := ""
		expectedHeader := http.Header{}
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})

	t.Run("Se intenta eliminar un producto pero no esta autorizado", func(t *testing.T) {
		// Arrange
		data := map[int]internal.Product{
			1: {
				ID:          1,
				Name:        "Producto 1",
				Quantity:    10,
				CodeValue:   "123456",
				IsPublished: true,
				Expiration:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
				Price:       100,
			},
		}
		repo := repository.NewRepositoryMap(data)
		service := service.NewProductService(repo)
		handler := handler.NewProductHandler(service)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/products/1", nil)

		// add id path param to the request
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		context := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
		req = req.WithContext(context)

		// Act
		handler.DeleteProduct(res, req)

		// Assert
		expectedCode := http.StatusUnauthorized
		expectedBody := `{"message":"Unauthorized", "status": 401}`
		expectedHeader := http.Header{
			"Content-Type": []string{"application/json"},
		}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})
}

// test using query params
func TestGetProductsByPriceGreaterThan(t *testing.T) {
	t.Run("Se solicitan los productos cuyo precio sea mayor a 100", func(t *testing.T) {
		// Arrange
		data := map[int]internal.Product{
			1: {
				ID:          1,
				Name:        "Producto 1",
				Quantity:    10,
				CodeValue:   "123456",
				IsPublished: true,
				Expiration:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
				Price:       100,
			},
			2: {
				ID:          2,
				Name:        "Producto 1",
				Quantity:    10,
				CodeValue:   "123456",
				IsPublished: true,
				Expiration:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
				Price:       110,
			},
		}
		repo := repository.NewRepositoryMap(data)
		service := service.NewProductService(repo)
		handler := handler.NewProductHandler(service)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/products?priceGt=100", nil)

		// Act
		handler.GetProductsByPriceGreaterThan(res, req)

		// Assert
		expectedCode := http.StatusOK
		expectedBody := `[{"id":2,"name":"Producto 1","quantity":10,"code_value":"123456","is_published":true,"expiration":"31/12/2021","price":110}]`
		expectedHeader := http.Header{
			"Content-Type": []string{"application/json"},
		}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})
}
