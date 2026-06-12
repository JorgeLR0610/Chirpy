package main

import (
	"log"
	"net/http"
	"time"

	"github.com/JorgeLR0610/Chirpy/internal/auth"
)

func(cfg *apiConfig) handlerRefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	storedToken, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil || storedToken.ExpiresAt.Before(time.Now().UTC()) || storedToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	newAccesToken, err := auth.MakeJWT(storedToken.UserID, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		log.Printf("There was an error creating refreshing a token access: %v", err)
		return
	}

	responseToken := struct {
		Token	string	`json:"token"`
	}{
		Token: newAccesToken,
	}

	respondWithJSON(w, http.StatusOK, responseToken)
}