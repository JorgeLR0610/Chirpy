package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func(cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpUUID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.DB.GetChirp(r.Context(), chirpUUID)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
		}

		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		log.Printf("There was an error retrieving the chirps: %v", err)
		return
	}
	
	resp := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, resp)
}