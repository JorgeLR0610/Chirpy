package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/JorgeLR0610/Chirpy/internal/auth"
	"github.com/JorgeLR0610/Chirpy/internal/database"
)


func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	type chirpParameters struct {
		Body string `json:"body"`	
	}

	decoder := json.NewDecoder(r.Body)
	chirpParams := chirpParameters{}
	if err := decoder.Decode(&chirpParams); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
	}

	userToken, err := auth.GetBearerToken(r.Header) 
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userUUID, err := auth.ValidateJWT(userToken, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	
	if len(chirpParams.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := replaceProfane(chirpParams.Body)

	newChirp, err := cfg.DB.CreateChirp(r.Context(), 
		database.CreateChirpParams{
			Body: cleanedBody,
			UserID: userUUID,
		})
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		log.Printf("There was an error creating the post: %v", err)
		return
	}

	post := Chirp {
		ID: newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body: newChirp.Body,
		UserID: newChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, post)
}


func replaceProfane(body string) string {

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}

	splitedBody := strings.Split(body, " ")

	for i, word := range splitedBody {
		if _, exists := badWords[strings.ToLower(word)]; exists {
			splitedBody[i] = "****"
		}
	}

	return strings.Join(splitedBody, " ")
}



