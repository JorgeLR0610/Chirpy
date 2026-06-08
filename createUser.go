package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func(cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type userEmail struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	usr := userEmail{}
	if err := decoder.Decode(&usr); err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		log.Printf("There was an error decoding the request: %v", err)
		return
	}

	newUser, err := cfg.DB.CreateUser(r.Context(), usr.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		log.Printf("There was an error creating the user: %v", err)
		return
	}

	user := User {
		ID: newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email: newUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)	
}