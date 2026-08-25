package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)


func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r * http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing given chirp id", err)
		return
	}

	chirpDB, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error retrieving chirp from database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID: chirpDB.ID,
		CreatedAt: chirpDB.CreatedAt,
		UpdatedAt: chirpDB.UpdatedAt,
		Body: chirpDB.Body,
		UserID: chirpDB.UserID,
	})
}