package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/JorgeLR0610/Chirpy/internal/auth"
	"github.com/JorgeLR0610/Chirpy/internal/database"
)

func(cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	usr := UserLoginParams{}
	if err := decoder.Decode(&usr); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
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

	// Create access token (1 hour) and refresh token
	userToken, err := auth.MakeJWT(u.ID, cfg.Secret)
	userRefreshToken := auth.MakeRefreshToken()

	// Create refresh token in DB, expiring on 60 days
	_, err = cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID: u.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})

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
		RefreshToken: userRefreshToken,
	}

	respondWithJSON(w, http.StatusOK, user)
}