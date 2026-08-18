package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Error string `json:"error"`
		Valid bool `json:"valid"`
	}

	params := parameters{}
	respBody := returnVals{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respBody.Error = fmt.Sprintf("Error decoding parameters: %s", err)
		respBody.Valid = false
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	

	if len(params.Body) > 140 {
		log.Printf("Chirp is too long")
		w.WriteHeader(400)
		return
	}
}
