package main

import (
	"fmt"
	"net/http"
	"path/filepath"
)

func main() {
	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	basePath := filepath.Dir(".")
	serveMux.Handle("/", http.FileServer(http.Dir(basePath)))

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error listening to server: %v\n", err)
	}
}
