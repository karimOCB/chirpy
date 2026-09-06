package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/karimOCB/chirpy/internal/auth"
)

func (cfg *apiConfig) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
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
		return
	}

	if body.ExpiresInSeconds == 0 || body.ExpiresInSeconds > 3600 {
		body.ExpiresInSeconds = 3600
	}

	token, err := auth.MakeJWT(userDB.ID, cfg.tokenSecret, time.Duration(body.ExpiresInSeconds))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating token", err)
		return
	}

	type LoginResponse struct {
		User
		Token string `json:"token"`
	}

	user := User{
		ID:        userDB.ID,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
		Email:     userDB.Email,
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		User:  user,
		Token: token,
	})
}
