package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/JorgeLR0610/Chirpy/internal/auth"
)

func(cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	usr := UserCredentials{}
	if err := decoder.Decode(&usr); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
	}

	u, err := cfg.DB.GetUser(r.Context(), usr.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		log.Printf("There was an error logging the user: %v", err)
		return
	}

	// Verify hash
	match, err := auth.CheckPasswordHash(usr.Password, u.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		log.Printf("There was an error while verifying password: %v", err)
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	user := User {
		ID: u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email: u.Email,
	}

	respondWithJSON(w, http.StatusOK, user)
}