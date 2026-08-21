package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

	cleanedBody := validateBadWords(params.Body)

	type returnVals struct {
		CleanedBody string  `json:"cleaned_body"`
	}
	
	payload := returnVals{
		CleanedBody: cleanedBody,
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


func validateBadWords(s string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}

	words := strings.Split(s, " ")

	for i, word := range(words) {
		lowWord := strings.ToLower(word)
		if _, ok := badWords[lowWord]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}