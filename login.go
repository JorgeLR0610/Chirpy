package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/JorgeLR0610/Chirpy/internal/auth"
)

func(cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	usr := UserLoginParams{}
	if err := decoder.Decode(&usr); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
	}

	if usr.ExpiresInSeconds <= 0 || usr.ExpiresInSeconds > 3600 {
		usr.ExpiresInSeconds = 3600
	}

	u, err := cfg.DB.GetUser(r.Context(), usr.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			auth.SimulatePasswordCheck(usr.Password)
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		log.Printf("There was an error finding the user email: %v", err)
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

	userToken, err := auth.MakeJWT(u.ID, cfg.Secret, time.Duration(usr.ExpiresInSeconds))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an error on our side")
		log.Printf("There was an error creating the JWT for user with ID %s: %v", u.ID, err)
		return
	}

	user := User {
		ID: u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email: u.Email,
		Token: userToken,
	}

	respondWithJSON(w, http.StatusOK, user)
}