package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/iankencruz/webbuilder/internal/repository"
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

func (s *AuthService) FindOrCreateUser(ctx context.Context, sub, provider, email, name, avatarURL string) (repository.User, error) {
	_, err := s.repo.GetUserBySub(ctx, sub)
	if err == nil {
		return s.repo.UpdateUser(ctx, repository.UpdateUserParams{
			Sub:       sub,
			Email:     email,
			Name:      pgtype.Text{String: name, Valid: name != ""},
			AvatarUrl: pgtype.Text{String: avatarURL, Valid: avatarURL != ""},
		})
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return repository.User{}, fmt.Errorf("looking up user by sub: %w", err)
	}

	return s.repo.CreateUser(ctx, repository.CreateUserParams{
		Sub:       sub,
		Provider:  provider,
		Email:     email,
		Name:      pgtype.Text{String: name, Valid: name != ""},
		AvatarUrl: pgtype.Text{String: avatarURL, Valid: avatarURL != ""},
	})
}

func (s *AuthService) GetByID(ctx context.Context, id int64) (repository.User, error) {
	return s.repo.GetUserByID(ctx, id)
}
