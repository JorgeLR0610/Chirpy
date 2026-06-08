package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"github.com/JorgeLR0610/Chirpy/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
	}
	
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := replaceProfane(params.Body)

	newChirp, err := cfg.DB.CreateChirp(r.Context(), 
		database.CreateChirpParams{
			Body: cleanedBody,
			UserID: params.UserID,
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



