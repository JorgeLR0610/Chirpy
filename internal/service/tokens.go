package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/JorgeLR0610/Chirpy/internal/auth"
	"github.com/JorgeLR0610/Chirpy/internal/repository"
)

var ErrTokenExpOrRev = errors.New("token expired or revoked")

type TokenService struct {
	repo *repository.Queries
}

func NewTokenService(repo *repository.Queries) *TokenService {
	return &TokenService{repo: repo}
}

func (s *TokenService) RefreshAccessToken(ctx context.Context, refreshToken, secret string) (string, error) {
	storedToken, err := s.repo.GetUserFromRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Printf("There was an error retrieving a refresh token: %v", err)
		return "", err
	}

	if storedToken.ExpiresAt.Before(time.Now().UTC()) || storedToken.RevokedAt.Valid {
		return "", ErrTokenExpOrRev
	}

	newAccesToken, err := auth.MakeJWT(storedToken.UserID, secret)
	if err != nil {
		log.Printf("There was an error refreshing a token access: %v", err)
		return "", err
	}

	return newAccesToken, nil
}

func (s *TokenService) RevokeRefreshToken(ctx context.Context, refreshTokenToRevoke string) error {
	err := s.repo.RevokeRefreshToken(ctx, repository.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time: time.Now().UTC(),
			Valid: true,
		},
		Token: refreshTokenToRevoke,
	})
	if err != nil {
		return err
	}

	return nil
}