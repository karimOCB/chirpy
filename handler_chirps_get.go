package main

import (
	"net/http"
)


func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r * http.Request) {
	chirpsDB, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving users", err)
		return
	}
	
	chirps := []Chirp{}
	for _, chirpDB := range(chirpsDB) {
		chirps = append(chirps, Chirp{
			ID: chirpDB.ID,
			CreatedAt: chirpDB.CreatedAt,
			UpdatedAt: chirpDB.UpdatedAt,
			Body: chirpDB.Body,
			UserID: chirpDB.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}