package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"

	mux := http.NewServeMux() 

	server := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}

	filepathRoot := http.Dir(".")
	mux.Handle("/", http.FileServer(filepathRoot))

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't listen to the server: %v", err)
	}
}

