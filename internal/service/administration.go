package service

import (
	"context"
	"log"

	"github.com/JorgeLR0610/Chirpy/internal/repository"
)

type AdminService struct {
	repo	*repository.Queries
}

func NewAdminService(repo *repository.Queries) *AdminService {
	return &AdminService{repo: repo}
}


func (s *AdminService) DeleteUsers(ctx context.Context) error {
	if err := s.repo.DeleteUsers(ctx); err != nil {
		log.Printf("There was an error truncating the table: %v", err)
		return err
	}

	return nil
}