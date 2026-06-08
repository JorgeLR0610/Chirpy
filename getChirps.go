package main

import (
	"log"
	"net/http"
)


func(cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.DB.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		log.Printf("There was an error retrieving the chirps: %v", err)
		return
	}

	resp := make([]Chirp, 0, len(chirps))

	for _, c := range chirps {
		resp = append(resp, Chirp{
			ID: c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body: c.Body,
			UserID: c.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, resp)
}