package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/JorgeLR0610/Chirpy/internal/auth"
	"github.com/JorgeLR0610/Chirpy/internal/repository"
	"github.com/google/uuid"
)

var ErrInvalidPassword = errors.New("incorrect email or password")

type LoginResult struct {
	User         repository.User
	AccessToken  string
	RefreshToken string
}

type UserService struct {
	repo	*repository.Queries
}

func NewUserService(repo *repository.Queries) *UserService {
	return &UserService{repo: repo}
}

func(s *UserService) CreateUser(ctx context.Context, email, password string) (repository.CreateUserRow, error) {

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		log.Printf("Hashing failed: %v", err)
		return repository.CreateUserRow{}, err
	}

	user, err := s.repo.CreateUser(ctx, repository.CreateUserParams{
		Email: email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("There was an error creating the user: %v", err)
		return repository.CreateUserRow{}, err
	}

	return user, nil
}

func(s *UserService) LoginUser(ctx context.Context, email, password, secret string) (LoginResult, error) {

user, err := s.repo.GetUser(ctx, email)
if err != nil {
	if errors.Is(err, sql.ErrNoRows) {
		auth.SimulatePasswordCheck(password)
		return LoginResult{}, ErrInvalidPassword
	}

	log.Printf("There was an error finding the user email")
	return LoginResult{}, err
}

	// Verify password
	match, err := auth.CheckPasswordHash(password, user.HashedPassword)
	if err != nil {
		log.Printf("There was an error while verifying password: %v", err)
		return LoginResult{}, err
	}

	if !match {
		return LoginResult{}, errors.New("Incorrect email or password")
	}

	accessToken, err := auth.MakeJWT(user.ID, secret)
	if err != nil {
		log.Printf("There was an error signing the access token: %v", err)
	}

	RefreshToken := auth.MakeRefreshToken()

	s.repo.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		Token: RefreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	return LoginResult{
		User: user,
		AccessToken: accessToken,
		RefreshToken: RefreshToken,
	}, nil

}

func(s *UserService) UpdateCredentials(ctx context.Context, email, newPassword string, userID uuid.UUID) (repository.UpdateEmailAndPasswordRow, error) {
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		log.Printf("Hashing failed: %v", err)
		return repository.UpdateEmailAndPasswordRow{}, err
	}

	newUser, err := s.repo.UpdateEmailAndPassword(ctx, repository.UpdateEmailAndPasswordParams{
		Email: email,
		HashedPassword: hashedPassword,
		UpdatedAt: time.Now().UTC(),
		ID: userID,
	})
	if err != nil {
		log.Printf("There was an error updating a user credentials: %v", err)
		return repository.UpdateEmailAndPasswordRow{}, err
	}

	return newUser, nil
}

