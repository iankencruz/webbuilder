package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/iankencruz/webbuilder/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	repo *repository.Queries
}

func NewAuthService(repo *repository.Queries) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (repository.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return repository.User{}, err
	}

	return s.repo.CreateUser(ctx, repository.CreateUserParams{
		Email:        strings.TrimSpace(strings.ToLower(email)),
		PasswordHash: string(hash),
	})
}

func (s *AuthService) Login(ctx context.Context, email, password string) (repository.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, strings.TrimSpace(strings.ToLower(email)))
	if err != nil {
		return repository.User{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return repository.User{}, ErrInvalidCredentials
	}

	return user, nil
}

func (s *AuthService) GetByID(ctx context.Context, id int64) (repository.User, error) {
	return s.repo.GetUserByID(ctx, id)
}
