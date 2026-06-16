package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/JorgeLR0610/Chirpy/internal/repository"
	"github.com/google/uuid"
)

var ErrChirpTooLong = errors.New("chirp is too long")
var ErrChirpTooShort = errors.New("chirp is too short")
var ErrNoRows = errors.New("resource not found")

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

	if len(body) < 1 {
		return repository.Chirp{}, ErrChirpTooShort
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
	chirps, err := s.repo.GetChirps(ctx)
    if err != nil {
		log.Printf("There was an error retrieving the chirps: %v", err)
		return nil, err
	}

	return chirps, nil
}

func (s *ChirpService) GetChirpsFromAuthor(ctx context.Context, UserID uuid.UUID) ([]repository.Chirp, error) {
	_, err := s.repo.GetUserByID(ctx, UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRows
		}
	}

	chirps, err := s.repo.GetChirpsFromAuthor(ctx, UserID)
    if err != nil {
		log.Printf("There was an error retrieving the chirps: %v", err)
		return nil, err
	}

	return chirps, nil
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

func (s *ChirpService) DeleteChirp(ctx context.Context, chirpID uuid.UUID) (error) {
	rows, err := s.repo.DeleteChirp(ctx, chirpID) 

	if err != nil {
		log.Printf("Could not delete a row in table chirps: %v", err)
		return err
	}

	if rows == 0 {
		return ErrNoRows
	}

	return nil

}