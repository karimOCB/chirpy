package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		respondWithError(w, 500, "Error decoding parameters")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}
	
	payload := returnVals{
		Valid: true,
	}
	
	respondWithJSON(w, 200, payload)
}


func respondWithError(w http.ResponseWriter, code int, msg string, /*passedErr error*/) {
	type returnVals struct {
		Error string `json:"error"`
	}

	respBody := returnVals{}	
	respBody.Error = msg

	respondWithJSON(w, code, respBody)	
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}