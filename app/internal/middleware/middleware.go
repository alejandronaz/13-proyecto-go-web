package middleware

import (
	"fmt"
	appHandler "goweb/app/internal/handler"
	"net/http"
	"time"

	"github.com/bootcamp-go/web/response"
)

func Auth(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader != "1234" {
			response.JSON(w, http.StatusUnauthorized, appHandler.ErrorResponse{
				Message: "Unauthorized",
				Status:  http.StatusUnauthorized,
			})
			return
		}

		// call the handler
		handler.ServeHTTP(w, r)
	})
}

func Logs(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		method := r.Method
		path := r.URL.Path
		timeReq := time.Now()
		tamano := r.ContentLength

		// call the handler
		handler.ServeHTTP(w, r)

		totalTime := time.Since(timeReq).Seconds()
		fmt.Printf("Method: %s - Path: %s - Time: %s - Tamaño: %db, Total time execution: %fs\n", method, path, timeReq.Format("02/01/2006 15:04:05"), tamano, totalTime)
	})
}
