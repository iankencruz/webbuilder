package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/iankencruz/webbuilder/internal/database/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// -- mocks
type mockAuthRepo struct {
	getUserByID  func(ctx context.Context, id int64) (repository.User, error)
	getUserBySub func(ctx context.Context, sub string) (repository.User, error)
	createUser   func(ctx context.Context, arg repository.CreateUserParams) (repository.User, error)
	updateUser   func(ctx context.Context, arg repository.UpdateUserParams) (repository.User, error)
}

func (m *mockAuthRepo) GetUserByID(ctx context.Context, id int64) (repository.User, error) {
	return m.getUserByID(ctx, id)
}

func (m *mockAuthRepo) GetUserBySub(ctx context.Context, sub string) (repository.User, error) {
	return m.getUserBySub(ctx, sub)
}

func (m *mockAuthRepo) CreateUser(ctx context.Context, arg repository.CreateUserParams) (repository.User, error) {
	return m.createUser(ctx, arg)
}

func (m *mockAuthRepo) UpdateUser(ctx context.Context, arg repository.UpdateUserParams) (repository.User, error) {
	return m.updateUser(ctx, arg)
}

// -- Test --

func TestFindOrCreateUser(t *testing.T) {
	ctx := context.Background()
	logger := slog.Default()

	existingUser := repository.User{
		ID:        1,
		Sub:       "existing-sub",
		Provider:  "google",
		Email:     "existing@example.com",
		FirstName: pgtype.Text{String: "Existing", Valid: true},
		LastName:  pgtype.Text{String: "User", Valid: true},
		AvatarUrl: pgtype.Text{String: "http://example.com/avatar.png", Valid: true},
	}

	tests := []struct {
		name      string
		sub       string
		provider  string
		email     string
		firstName string
		lastName  string
		avatarURL string
		repo      *mockAuthRepo
		wantUser  repository.User
		wantErr   bool
	}{
		{
			name:      "Update existing user info",
			sub:       "existing-sub",
			provider:  "google",
			email:     "existing@example.com",
			firstName: "Updated",
			lastName:  "User",
			avatarURL: "http://example.com/new-avatar.png",
			repo: &mockAuthRepo{
				getUserBySub: func(ctx context.Context, sub string) (repository.User, error) {
					return existingUser, nil
				},
				updateUser: func(ctx context.Context, arg repository.UpdateUserParams) (repository.User, error) {
					return repository.User{
						ID:        1,
						Sub:       "existing-sub",
						Provider:  "google",
						Email:     arg.Email,
						FirstName: pgtype.Text{String: arg.FirstName.String, Valid: true},
						LastName:  pgtype.Text{String: arg.LastName.String, Valid: true},
						AvatarUrl: pgtype.Text{String: arg.AvatarUrl.String, Valid: true},
					}, nil
				},
			},
			wantUser: repository.User{
				ID:        1,
				Sub:       "existing-sub",
				Provider:  "google",
				Email:     "existing@example.com",
				FirstName: pgtype.Text{String: "Updated", Valid: true},
				LastName:  pgtype.Text{String: "User", Valid: true},
				AvatarUrl: pgtype.Text{String: "http://example.com/new-avatar.png", Valid: true},
			},
			wantErr: false,
		},
		{
			name: "Create new user",
			repo: &mockAuthRepo{
				getUserBySub: func(ctx context.Context, sub string) (repository.User, error) {
					return repository.User{}, pgx.ErrNoRows
				},
				createUser: func(ctx context.Context, arg repository.CreateUserParams) (repository.User, error) {
					return repository.User{
						ID:        2,
						Sub:       arg.Sub,
						Provider:  arg.Provider,
						Email:     arg.Email,
						FirstName: pgtype.Text{String: arg.FirstName.String, Valid: true},
						LastName:  pgtype.Text{String: arg.LastName.String, Valid: true},
						AvatarUrl: pgtype.Text{String: arg.AvatarUrl.String, Valid: true},
					}, nil
				},
			},
			sub:       "new-sub",
			provider:  "google",
			email:     "new@example.com",
			firstName: "New",
			lastName:  "User",
			avatarURL: "http://example.com/new-avatar.png",
			wantUser: repository.User{
				ID:        2,
				Sub:       "new-sub",
				Provider:  "google",
				Email:     "new@example.com",
				FirstName: pgtype.Text{String: "New", Valid: true},
				LastName:  pgtype.Text{String: "User", Valid: true},
				AvatarUrl: pgtype.Text{String: "http://example.com/new-avatar.png", Valid: true},
			},
			wantErr: false,
		},
		{
			name: "Unexpected repo error",
			repo: &mockAuthRepo{
				getUserBySub: func(ctx context.Context, sub string) (repository.User, error) {
					return repository.User{}, pgx.ErrNoRows
				},
				createUser: func(ctx context.Context, arg repository.CreateUserParams) (repository.User, error) {
					return repository.User{}, pgx.ErrTxClosed
				},
			},
			sub:      "new-sub",
			provider: "google",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAuthService(logger, tt.repo)

			got, err := service.FindOrCreateUser(ctx, tt.sub, tt.provider, tt.email, tt.firstName, tt.lastName, tt.avatarURL)

			if (err != nil) != tt.wantErr {
				t.Fatalf("FindOrCreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if got.ID != tt.wantUser.ID ||
					got.Sub != tt.wantUser.Sub ||
					got.Provider != tt.wantUser.Provider ||
					got.Email != tt.wantUser.Email ||
					got.FirstName.String != tt.wantUser.FirstName.String ||
					got.LastName.String != tt.wantUser.LastName.String ||
					got.AvatarUrl.String != tt.wantUser.AvatarUrl.String {
					t.Errorf("FindOrCreateUser() got = %+v, want %+v", got, tt.wantUser)
				}
				if got.Email != tt.wantUser.Email {
					t.Errorf("FindOrCreateUser() got email = %s, want %s", got.Email, tt.wantUser.Email)
				}
				if got.FirstName.String != tt.wantUser.FirstName.String {
					t.Errorf("FindOrCreateUser() got firstName = %s, want %s", got.FirstName.String, tt.wantUser.FirstName.String)
				}
				if got.LastName.String != tt.wantUser.LastName.String {
					t.Errorf("FindOrCreateUser() got lastName = %s, want %s", got.LastName.String, tt.wantUser.LastName.String)
				}
				if got.AvatarUrl.String != tt.wantUser.AvatarUrl.String {
					t.Errorf("FindOrCreateUser() got avatarURL = %s, want %s", got.AvatarUrl.String, tt.wantUser.AvatarUrl.String)
				}
				if got.Sub != tt.wantUser.Sub {
					t.Errorf("FindOrCreateUser() got sub = %s, want %s", got.Sub, tt.wantUser.Sub)
				}
				if got.Provider != tt.wantUser.Provider {
					t.Errorf("FindOrCreateUser() got provider = %s, want %s", got.Provider, tt.wantUser.Provider)
				}
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	logger := slog.Default()

	tests := []struct {
		name     string
		repo     *mockAuthRepo
		id       int64
		wantUser repository.User
		wantErr  bool
	}{
		{
			name: "User found",
			repo: &mockAuthRepo{
				getUserByID: func(ctx context.Context, id int64) (repository.User, error) {
					return repository.User{
						ID:       id,
						Sub:      "test-sub",
						Provider: "google",
						Email:    "found@example.com",
					}, nil
				},
			},
			id: 1,
			wantUser: repository.User{
				ID:       1,
				Sub:      "test-sub",
				Provider: "google",
				Email:    "found@example.com",
			},
			wantErr: false,
		},
		{
			name: "User not found",
			repo: &mockAuthRepo{
				getUserByID: func(ctx context.Context, id int64) (repository.User, error) {
					return repository.User{}, pgx.ErrNoRows
				},
			},
			id:      99,
			wantErr: true,
		},
		{
			name: "Unexpected repo error",
			repo: &mockAuthRepo{
				getUserByID: func(ctx context.Context, id int64) (repository.User, error) {
					return repository.User{}, pgx.ErrTxClosed
				},
			},
			id:      1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAuthService(logger, tt.repo)

			got, err := service.GetByID(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr %v, got error %v", tt.wantErr, err)
			}

			if !tt.wantErr {
				if got.ID != tt.wantUser.ID ||
					got.Sub != tt.wantUser.Sub ||
					got.Provider != tt.wantUser.Provider ||
					got.Email != tt.wantUser.Email {
					t.Errorf("GetByID() got = %+v, want %+v", got, tt.wantUser)
				}
				if got.Email != tt.wantUser.Email {
					t.Errorf("Email: want %s, got %s", tt.wantUser.Email, got.Email)
				}
				if got.Sub != tt.wantUser.Sub {
					t.Errorf("Sub: want %s, got %s", tt.wantUser.Sub, got.Sub)
				}
				if got.Provider != tt.wantUser.Provider {
					t.Errorf("Provider: want %s, got %s", tt.wantUser.Provider, got.Provider)
				}
			}
		})
	}
}
