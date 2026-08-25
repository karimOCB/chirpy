package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/karimOCB/chirpy/internal/auth"
)


func (cfg *apiConfig) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	body := parameters{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request body", err)
	}

	userDB, err := cfg.dbQueries.GetUserByEmail(r.Context(), body.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User does not exist in the database -> 401 Unauthorized
			respondWithError(w, http.StatusUnauthorized, "Incorrect email", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error retrieving user from database", err)
		return
	}

	ok, err := auth.CheckPasswordHash(body.Password, userDB.HashedPassword)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error checking user password", nil)
		return
	}
	
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password", nil)
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: userDB.ID,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
		Email: userDB.Email,
	})
}