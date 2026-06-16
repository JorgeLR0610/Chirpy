package api

import (
	"time"

	"github.com/google/uuid"
)

type userLoginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsChirpyRed	 bool		`json:"is_chirpy_red"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type userCreationResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed	 bool		`json:"is_chirpy_red"`
}

type chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type userRegisterParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type userLoginParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type userNewCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userCredentialsUpdateResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed	 bool	`json:"is_chirpy_red"`
}

type userSubscriptionUpgrade struct {
	Event	string	 `json:"event"`
	Data	UserUpgradeData `json:"data"`
}

type UserUpgradeData struct {
	UserID	string	`json:"user_id"`
}