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
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(basePath))))

	serveMux.HandleFunc(
		"/healthz/",
		func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
			writer.WriteHeader(200)
			writer.Write([]byte("OK"))
		},
	)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error listening to server: %v\n", err)
	}
}
