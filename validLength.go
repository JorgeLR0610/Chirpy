package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)


func handlerBody(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters")
		return
	}
	
	// params is a struct with data populated successfully
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Call the profaneValidation
	cleanedBody := replaceProfane(params.Body)

	type validResponse struct {
		Cleanedbody string `json:"cleaned_body"`
	}

	respondWithJSON(w, http.StatusOK, validResponse{
		Cleanedbody: cleanedBody,
	})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {

	w.Header().Set("Content-Type", "application/json")

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
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



