package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/JorgeLR0610/Chirpy/internal/auth"
	"github.com/JorgeLR0610/Chirpy/internal/database"
)

func(cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	err = cfg.DB.RevokeRefreshToken(r.Context(), database.			RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time: time.Now().UTC(),
			Valid: true,
		},
		Token: refreshToken,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return		
	}

	w.WriteHeader(http.StatusNoContent)
}