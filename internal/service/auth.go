package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/iankencruz/webbuilder/internal/database/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthRepository interface {
	GetUserByID(ctx context.Context, id int64) (repository.User, error)
	GetUserBySub(ctx context.Context, sub string) (repository.User, error)
	CreateUser(ctx context.Context, arg repository.CreateUserParams) (repository.User, error)
	UpdateUser(ctx context.Context, arg repository.UpdateUserParams) (repository.User, error)
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) FindOrCreateUser(ctx context.Context, sub, provider, email, firstName, lastName, avatarURL string) (repository.User, error) {
	log.Printf("DEBUG FindOrCreateUser sub=%s provider=%s email=%s firstName=%s lastname=%s", sub, provider, email, firstName, lastName)

	_, err := s.repo.GetUserBySub(ctx, sub)
	if err == nil {
		return s.repo.UpdateUser(ctx, repository.UpdateUserParams{
			Sub:       sub,
			Email:     email,
			FirstName: pgtype.Text{String: firstName, Valid: firstName != ""},
			LastName:  pgtype.Text{String: lastName, Valid: lastName != ""},
			AvatarUrl: pgtype.Text{String: avatarURL, Valid: avatarURL != ""},
		})
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		log.Printf("DEBUG unexpected repository error: %v", err)
		return repository.User{}, fmt.Errorf("looking up user by sub: %w", err)
	}
	log.Printf("DEBUG creating new user")
	return s.repo.CreateUser(ctx, repository.CreateUserParams{
		Sub:       sub,
		Provider:  provider,
		Email:     email,
		FirstName: pgtype.Text{String: firstName, Valid: firstName != ""},
		LastName:  pgtype.Text{String: lastName, Valid: lastName != ""},
		AvatarUrl: pgtype.Text{String: avatarURL, Valid: avatarURL != ""},
	})
}

func (s *AuthService) GetByID(ctx context.Context, id int64) (repository.User, error) {
	return s.repo.GetUserByID(ctx, id)
}
