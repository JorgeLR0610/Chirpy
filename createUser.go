package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func(cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type userEmail struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	usr := userEmail{}
	if err := decoder.Decode(&usr); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid fields")
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