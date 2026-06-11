package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/JorgeLR0610/Chirpy/internal/auth"
	"github.com/JorgeLR0610/Chirpy/internal/database"
)

func(cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	usr := UserRegisterParams{}
	if err := decoder.Decode(&usr); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
	}

	hashedPassword, err := auth.HashPassword(usr.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		log.Printf("Hashing failed: %v", err)
		return
	}

	newUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email: usr.Email,
		HashedPassword: hashedPassword,
	})

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