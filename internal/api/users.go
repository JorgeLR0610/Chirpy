package api

import (
	"encoding/json"
	"net/http"

	"github.com/JorgeLR0610/Chirpy/internal/auth"
	"github.com/JorgeLR0610/Chirpy/internal/service"
)

type UserHandler struct {
	service *service.UserService
	secret  string
}

func NewUserHandler(svc *service.UserService, secret string) *UserHandler {
	return &UserHandler{service: svc, secret: secret}
}

func (h *UserHandler) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var user userRegisterParams
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
	}

	newUser, err := h.service.CreateUser(r.Context(), user.Email, user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		return
	}

	respondWithJSON(w, http.StatusCreated, userCreationResponse{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	})
}

func (h *UserHandler) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials userLoginParams
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
	}

	result, err := h.service.LoginUser(r.Context(), credentials.Email, credentials.Password, h.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	respondWithJSON(w, http.StatusOK, userLoginResponse{
		ID:           result.User.ID,
		CreatedAt:    result.User.CreatedAt,
		UpdatedAt:    result.User.UpdatedAt,
		Email:        result.User.Email,
		Token:        result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}

func (h *UserHandler) HandlerUpdateCredentials(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(accessToken, h.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var newCredentials userNewCredentials
	if err := json.NewDecoder(r.Body).Decode(&newCredentials); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
	}

	updatedUser, err := h.service.UpdateCredentials(r.Context(), newCredentials.Email, newCredentials.Password, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		return
	}

	respondWithJSON(w, http.StatusOK, userCredentialsUpdateResponse{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	})

}
