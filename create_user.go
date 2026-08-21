package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	
	ctx := r.Context()
	body := parameters{}

	err := json.NewDecoder(r.Body).Decode(&body)
	
	if err != nil {
		respondWithError(w, 500, "Error decoding parameters")
		return
	}

	db_user, err := cfg.dbQueries.CreateUser(ctx, body.Email)
	if err != nil {
		respondWithError(w, 500, "Error creating user")
		return
	}
	
	user := User{}
	user.ID = db_user.ID
	user.CreatedAt = db_user.CreatedAt
	user.UpdatedAt = db_user.UpdatedAt
	user.Email = db_user.Email

	respondWithJSON(w, 201, user)
}