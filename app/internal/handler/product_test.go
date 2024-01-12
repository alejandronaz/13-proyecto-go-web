package handler_test

import (
	"goweb/app/internal"
	"goweb/app/internal/handler"
	"goweb/app/internal/repository"
	"goweb/app/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
