package service

import (
	"context"
	"errors"
	"strings"

	"github.com/JorgeLR0610/Chirpy/internal/repository"
	"github.com/google/uuid"
)

var ErrChirpTooLong = errors.New("Chirp is too long")

type ChirpService struct {
	repo	*repository.Queries
}

func NewChirpService(repo *repository.Queries) *ChirpService {
    return &ChirpService{repo: repo}
}

func (s *ChirpService) CreateChirp(ctx context.Context, body string, userID uuid.UUID) (repository.Chirp, error) {
	    if len(body) > 140 {
        return repository.Chirp{}, ErrChirpTooLong
    }

    cleaned := replaceProfane(body)

    return s.repo.CreateChirp(ctx, repository.CreateChirpParams{
        Body:   cleaned,
        UserID: userID,
    })
}

func (s *ChirpService) GetChirp(ctx context.Context, id uuid.UUID) (repository.Chirp, error) {
    return s.repo.GetChirp(ctx, id)
}

func (s *ChirpService) GetChirps(ctx context.Context) ([]repository.Chirp, error) {
    return s.repo.GetChirps(ctx)
}

func replaceProfane(body string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}

	words := strings.Split(body, " ")

	for i, word := range words {
		if _, exists := badWords[strings.ToLower(word)]; exists {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}