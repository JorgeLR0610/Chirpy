package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/JorgeLR0610/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Post struct {
	ID			uuid.UUID	`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Body		string		`json:"body"`
	UserID		uuid.UUID	`json:"user_id"`
}


func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters")
		return
	}
	
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := replaceProfane(params.Body)

	newPost, err := cfg.DB.CreatePost(r.Context(), 
		database.CreatePostParams{
			Body: cleanedBody,
			UserID: params.UserID,
		})
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		log.Printf("There was an error creating the post: %v", err)
		return
	}

	post := Post {
		ID: newPost.ID,
		CreatedAt: newPost.CreatedAt,
		UpdatedAt: newPost.UpdatedAt,
		Body: newPost.Body,
		UserID: newPost.UserID,
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



