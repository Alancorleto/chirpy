package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) metricsResponse(writer http.ResponseWriter, request *http.Request) {
	htmlResponse := `
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`
	writer.Header().Add("Content-Type", "text/html")
	writer.WriteHeader(200)
	writer.Write([]byte(fmt.Sprintf(htmlResponse, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(writer, request)
		},
	)
}
