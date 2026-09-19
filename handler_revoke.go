package main

import (
	"net/http"

	"github.com/karimOCB/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting bearer token", err)
		return
	}

	execResult, err := cfg.dbQueries.UpdateRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting token from database", err)
		return
	}

	rowsAffected, err := execResult.RowsAffected()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error checking affected rows", err)
		return
	}

	if rowsAffected == 0 {
		respondWithError(w, http.StatusUnauthorized, "error getting bearer token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
