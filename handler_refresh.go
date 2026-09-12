package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/karimOCB/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting bearer token", err)
		return
	}

	dbRefreshToken, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Refresh Token does not exist in the database -> 401 Unauthorized
			respondWithError(w, http.StatusUnauthorized, "incorrect refresh token", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error getting token from database", err)
		return
	}

	if time.Now().After(dbRefreshToken.ExpiresAt) || dbRefreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token expired or revoked", nil)
		return
	}

	accessToken, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.tokenSecret, time.Hour*1)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating token", err)
		return
	}

	type payload struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, payload{
		Token: accessToken,
	})
}
