package api

import (
	"errors"
	"net/http"

	"github.com/JorgeLR0610/Chirpy/internal/auth"
	"github.com/JorgeLR0610/Chirpy/internal/service"
)

type TokensHandler struct {
	service	service.TokenService
	secret	string
}

func NewRefreshAccessTokenHandler(tokenSvc *service.TokenService, secret string) *TokensHandler {
	return &TokensHandler{service: *tokenSvc, secret: secret}
}

func (h *TokensHandler) HandlerRefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	newAccessToken, err := h.service.RefreshAccessToken(r.Context(), refreshToken, h.secret)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpOrRev){
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		return
	}

	responseToken := struct {
		Token	string	`json:"token"`
	}{
		Token: newAccessToken,
	}

	respondWithJSON(w, http.StatusOK, responseToken)
}

func (h *TokensHandler) HandlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = h.service.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, service.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		
		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}