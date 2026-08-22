package main

import "net/http"

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if !cfg.platform {
		respondWithError(w, http.StatusForbidden, "cannot reset users information", nil)
	}

	cfg.dbQueries.DeleteUser(r.Context())

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}
