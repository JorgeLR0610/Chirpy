package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"

	"github.com/JorgeLR0610/Chirpy/internal/auth"
	"github.com/JorgeLR0610/Chirpy/internal/repository"
	"github.com/JorgeLR0610/Chirpy/internal/service"
	"github.com/google/uuid"
)

type ChirpHandler struct {
	service *service.ChirpService
	secret  string
}

func NewChirpHandler(svc *service.ChirpService, secret string) *ChirpHandler {
	return &ChirpHandler{service: svc, secret: secret}
}

func (h *ChirpHandler) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	type chirpParameters struct {
		Body string `json:"body"`
	}

	var chirpParams chirpParameters
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&chirpParams); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	requestToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(requestToken, h.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	newChirp, err := h.service.CreateChirp(r.Context(), chirpParams.Body, userID)
	if err != nil {
		if errors.Is(err, service.ErrChirpTooLong) || errors.Is(err, service.ErrChirpTooShort) {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	})
}

func (h *ChirpHandler) HandlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpUUID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	requestedChirp, err := h.service.GetChirp(r.Context(), chirpUUID)

	if err != nil {
		if errors.Is(err, service.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		log.Printf("There was an error retrieving the chirps: %v", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp{
		ID:        requestedChirp.ID,
		CreatedAt: requestedChirp.CreatedAt,
		UpdatedAt: requestedChirp.UpdatedAt,
		Body:      requestedChirp.Body,
		UserID:    requestedChirp.UserID,
	})
}

func (h *ChirpHandler) HandlerGetChirps(w http.ResponseWriter, r *http.Request) {

	sorting := r.URL.Query().Get("sort")

	var chirps []repository.Chirp
	var err error

	if authorIDStr := r.URL.Query().Get("author_id"); authorIDStr != "" {
		authorUUID, parseErr := uuid.Parse(authorIDStr)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
		chirps, err = h.service.GetChirpsFromAuthor(r.Context(), authorUUID)
	} else {
		chirps, err = h.service.GetChirps(r.Context())
	}
	
	if err != nil {
		if errors.Is(err, service.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Author not found")
			return
		}
		
		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		return
	}

	resp := make([]chirp, 0, len(chirps))

	for _, c := range chirps {
		resp = append(resp, chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	if sorting == "desc" {
		slices.SortFunc(resp, func(a, b chirp) int {
			return b.CreatedAt.Compare(a.CreatedAt)
	})
}
	respondWithJSON(w, http.StatusOK, resp)
}

func (h *ChirpHandler) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	requestUserID, err := auth.ValidateJWT(accessToken, h.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	chirpUUID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}	

	chirpToDelete, err := h.service.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		if errors.Is(err, service.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		return
	}

	if chirpToDelete.UserID != requestUserID {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	if err := h.service.DeleteChirp(r.Context(), chirpToDelete.ID); err != nil {
		if errors.Is(err, service.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Not found")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}