package api

import (
	"encoding/json"
	"net/http"

	"github.com/JorgeLR0610/Chirpy/internal/service"
)

type UserHandler struct {
	service *service.UserService
	secret string
}

func NewUserHandler(svc *service.UserService, secret string) *UserHandler {
	return &UserHandler{service: svc, secret: secret}
}

func (h *UserHandler) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var user UserRegisterParams
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
	}

	newUser, err := h.service.CreateUser(r.Context(), user.Email, user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was an error on our end")
		return
	}

	respondWithJSON(w, http.StatusCreated, UserCreationResponse{
		ID: newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email: newUser.Email,
	})
}

func (h *UserHandler) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials UserLoginParams
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid fields")
		return
	}

	result, err := h.service.LoginUser(r.Context(), credentials.Email, credentials.Password, h.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	respondWithJSON(w, http.StatusOK, UserLoginResponse{
		ID: result.User.ID,
		CreatedAt: result.User.CreatedAt,
		UpdatedAt: result.User.UpdatedAt,
		Email: result.User.Email,
		Token: result.AccessToken,
		RefreshToken: result.RefreshToken,
	})

}
