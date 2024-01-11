package main

import (
	"fmt"
	"goweb/app/internal/application"
)

func main() {

	// create the server
	server := application.NewServer(":8080")

	// run the server
	err := server.Start()
	if err != nil {
		fmt.Println(err)
	}

}
