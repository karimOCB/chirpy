package main

import (
	"strings"
)


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