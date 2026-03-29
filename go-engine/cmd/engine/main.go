package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	metrics "github.com/wendelvx/pi-devops.git/internal/metrics"
)

func main() {
	startedAt := time.Now()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		fmt.Fprintf(writer, `{"service":"go-engine","status":"pong","path":"%s"}`, request.URL.Path)
	})
	mux.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; version=0.0.4")
		writer.WriteHeader(http.StatusOK)
		fmt.Fprint(writer, metrics.Render(startedAt))
	})

	address := "0.0.0.0:" + port
	log.Printf("go-engine listening on %s", address)
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal(err)
	}
}
