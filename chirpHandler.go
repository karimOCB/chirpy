package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/karimOCB/chirpy/internal/database"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	} 

	params	:= parameters{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request body", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	params.Body = validateBadWords(params.Body)

	chirpDB, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{ Body: params.Body, UserID: params.UserID })
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating chirp", err)
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirpDB.ID,
		CreatedAt: chirpDB.CreatedAt,
		UpdatedAt: chirpDB.UpdatedAt,
		Body: chirpDB.Body,
		UserID: chirpDB.UserID,
	})
}