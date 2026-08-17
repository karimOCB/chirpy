package main

import (
	"fmt"
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
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/metrics", cfg.metricsHandler)
	mux.HandleFunc("/reset", cfg.resetHandler)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't listen to the server: %v", err)
	}
}


func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}


func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	fmt.Fprintf(w, "Hits: %d", cfg.fileserverHits.Load())
}



func (cfg *apiConfig) resetHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	fmt.Fprintf(w, "Reset. Hits: %d", cfg.fileserverHits.Load())
}


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}