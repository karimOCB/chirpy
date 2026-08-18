package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	port := "8080"
	cfg := apiConfig{}

	mux := http.NewServeMux() 

	server := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}

	filepathRoot := http.Dir(".")
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(filepathRoot))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't listen to the server: %v", err)
	}
}