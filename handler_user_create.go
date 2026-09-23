package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/karimOCB/chirpy/internal/auth"
	"github.com/karimOCB/chirpy/internal/database"
)

type User struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		IsChirpyRed bool `json:"is_chirpy_red"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	
	ctx := r.Context()
	body := parameters{}

	err := json.NewDecoder(r.Body).Decode(&body)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	hash, err := auth.HashPassword(body.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	db_user, err := cfg.dbQueries.CreateUser(ctx, database.CreateUserParams{
		Email: body.Email,
		HashedPassword: hash,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}
	
	user := User{
		ID: db_user.ID,
		CreatedAt: db_user.CreatedAt,
		UpdatedAt: db_user.UpdatedAt,
		Email: db_user.Email,
		IsChirpyRed: db_user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, user)
}