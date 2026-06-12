package pages

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iankencruz/webbuilder/internal/database/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v5"
)

// --- Mock ---

type mockRepository struct {
	createPage    func(ctx context.Context, arg repository.CreatePageParams) (repository.Page, error)
	getPageByID   func(ctx context.Context, id int64) (repository.Page, error)
	getPageBySlug func(ctx context.Context, slug string) (repository.Page, error)
	listPages     func(ctx context.Context) ([]repository.Page, error)
	updatePage    func(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error)
	deleteByID    func(ctx context.Context, id int64) error
	deleteBySlug  func(ctx context.Context, slug string) error
}

func (m *mockRepository) CreatePage(ctx context.Context, arg repository.CreatePageParams) (repository.Page, error) {
	return m.createPage(ctx, arg)
}

func (m *mockRepository) GetPageByID(ctx context.Context, id int64) (repository.Page, error) {
	return m.getPageByID(ctx, id)
}

func (m *mockRepository) GetPageBySlug(ctx context.Context, slug string) (repository.Page, error) {
	return m.getPageBySlug(ctx, slug)
}

func (m *mockRepository) ListPages(ctx context.Context) ([]repository.Page, error) {
	return m.listPages(ctx)
}

func (m *mockRepository) UpdatePage(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error) {
	return m.updatePage(ctx, arg)
}

func (m *mockRepository) DeleteByID(ctx context.Context, id int64) error {
	return m.deleteByID(ctx, id)
}

func (m *mockRepository) DeleteBySlug(ctx context.Context, slug string) error {
	return m.deleteBySlug(ctx, slug)
}

// --- Helpers ---

func newHandler() *Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewHandler(logger, &mockRepository{})
}

func newTestContext(e *echo.Echo, method, target, body string) (*echo.Context, *httptest.ResponseRecorder) {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, target, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

// --- Tests ---

func TestListPages(t *testing.T) {
	tests := []struct {
		name       string
		repo       mockRepository
		wantStatus int
		wantLen    int
	}{
		{
			name: "returns all pages",
			repo: mockRepository{
				listPages: func(ctx context.Context) ([]repository.Page, error) {
					return []repository.Page{
						{ID: 1, Title: "Home", Slug: "home"},
						{ID: 2, Title: "About", Slug: "about"},
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantLen:    2,
		},
		{
			name: "returns empty list",
			repo: mockRepository{
				listPages: func(ctx context.Context) ([]repository.Page, error) {
					return []repository.Page{}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantLen:    0,
		},
		{
			name: "db error returns 500",
			repo: mockRepository{
				listPages: func(ctx context.Context) ([]repository.Page, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			h := newHandler()
			c, rec := newTestContext(e, http.MethodGet, "/api/pages", "")

			if err := h.ListPages(c); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("status mismatch\n  want: %d\n  got:  %d\n  body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}

			if tt.wantStatus == http.StatusOK {
				var pages []repository.Page
				if err := json.Unmarshal(rec.Body.Bytes(), &pages); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(pages) != tt.wantLen {
					t.Errorf("length mismatch\n  want: %d\n  got:  %d", tt.wantLen, len(pages))
				}
			}
		})
	}
}

func TestCreatePage(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		repo       mockRepository
		wantStatus int
		wantTitle  string
		wantSlug   string
	}{
		{
			name: "valid request creates page",
			body: `{"title":"Contact","slug":"contact","status":"draft"}`,
			repo: mockRepository{
				createPage: func(ctx context.Context, arg repository.CreatePageParams) (repository.Page, error) {
					return repository.Page{ID: 3, Title: arg.Title, Slug: arg.Slug}, nil
				},
			},
			wantStatus: http.StatusCreated,
			wantTitle:  "Contact",
			wantSlug:   "contact",
		},
		{
			name:       "missing title returns 400",
			body:       `{"slug":"contact"}`,
			repo:       mockRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing slug returns 400",
			body:       `{"title":"Contact"}`,
			repo:       mockRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "db error returns 500",
			body: `{"title":"Contact","slug":"contact"}`,
			repo: mockRepository{
				createPage: func(ctx context.Context, arg repository.CreatePageParams) (repository.Page, error) {
					return repository.Page{}, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			h := newHandler()
			c, rec := newTestContext(e, http.MethodPost, "/api/pages", tt.body)

			if err := h.CreatePage(c); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("status mismatch\n  want: %d\n  got:  %d\n  body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}

			if tt.wantTitle != "" || tt.wantSlug != "" {
				var page repository.Page
				if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
					t.Fatalf("failed to decode response: %v\n  raw: %s", err, rec.Body.String())
				}
				if page.Title != tt.wantTitle {
					t.Errorf("title mismatch\n  want: %q\n  got:  %q", tt.wantTitle, page.Title)
				}
				if page.Slug != tt.wantSlug {
					t.Errorf("slug mismatch\n  want: %q\n  got:  %q", tt.wantSlug, page.Slug)
				}
			}
		})
	}
}

func TestGetPage(t *testing.T) {
	tests := []struct {
		name       string
		slug       string
		repo       mockRepository
		wantStatus int
	}{
		{
			name: "found returns 200",
			slug: "home",
			repo: mockRepository{
				getPageBySlug: func(ctx context.Context, slug string) (repository.Page, error) {
					return repository.Page{ID: 1, Title: "Home", Slug: slug}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not found returns 404",
			slug: "missing",
			repo: mockRepository{
				getPageBySlug: func(ctx context.Context, slug string) (repository.Page, error) {
					return repository.Page{}, pgx.ErrNoRows
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "db error returns 500",
			slug: "home",
			repo: mockRepository{
				getPageBySlug: func(ctx context.Context, slug string) (repository.Page, error) {
					return repository.Page{}, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			h := newHandler()
			c, rec := newTestContext(e, http.MethodGet, "/api/pages/"+tt.slug, "")
			c.SetPathValues(echo.PathValues{{Name: "slug", Value: tt.slug}})

			if err := h.GetPage(c); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("status mismatch\n  want: %d\n  got:  %d\n  body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestUpdatePage(t *testing.T) {
	tests := []struct {
		name       string
		slug       string
		body       string
		repo       mockRepository
		wantStatus int
		wantTitle  string
		wantSlug   string
	}{
		{
			name: "valid request updates page",
			slug: "home",
			body: `{"title":"Home Updated","status":"published"}`,
			repo: mockRepository{
				updatePage: func(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error) {
					return repository.Page{ID: 1, Title: arg.Title, Slug: arg.Slug}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantTitle:  "Home Updated",
			wantSlug:   "home",
		},
		{
			name:       "missing slug returns 400",
			slug:       "",
			body:       `{"title":"Home Updated"}`,
			repo:       mockRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not found returns 404",
			slug: "nonexistent",
			body: `{"title":"Home Updated","status":"published"}`,
			repo: mockRepository{
				updatePage: func(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error) {
					return repository.Page{}, pgx.ErrNoRows
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "db error returns 500",
			slug: "home",
			body: `{"title":"Home Updated","status":"published"}`,
			repo: mockRepository{
				updatePage: func(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error) {
					return repository.Page{}, errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			h := newHandler()

			req := httptest.NewRequest(http.MethodPut, "/api/pages/"+tt.slug, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.slug != "" {
				c.SetPathValues(echo.PathValues{{Name: "slug", Value: tt.slug}})
			}

			if err := h.UpdatePage(c); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("status mismatch\n  want: %d\n  got:  %d\n  body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}

			if tt.wantTitle != "" || tt.wantSlug != "" {
				var page repository.Page
				if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
					t.Fatalf("failed to decode response: %v\n  raw: %s", err, rec.Body.String())
				}
				if page.Title != tt.wantTitle {
					t.Errorf("title mismatch\n  want: %q\n  got:  %q", tt.wantTitle, page.Title)
				}
				if page.Slug != tt.wantSlug {
					t.Errorf("slug mismatch\n  want: %q\n  got:  %q", tt.wantSlug, page.Slug)
				}
			}
		})
	}
}

func TestDeletePage(t *testing.T) {
	tests := []struct {
		name       string
		slug       string
		repo       mockRepository
		wantStatus int
	}{
		{
			name: "valid slug deletes page",
			slug: "home",
			repo: mockRepository{
				deleteBySlug: func(ctx context.Context, slug string) error {
					return nil
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "not found returns 404",
			slug: "nonexistent",
			repo: mockRepository{
				deleteBySlug: func(ctx context.Context, slug string) error {
					return pgx.ErrNoRows
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "db error returns 500",
			slug: "home",
			repo: mockRepository{
				deleteBySlug: func(ctx context.Context, slug string) error {
					return errors.New("db error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			h := newHandler()

			req := httptest.NewRequest(http.MethodDelete, "/api/pages/"+tt.slug, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.slug != "" {
				c.SetPathValues(echo.PathValues{{Name: "slug", Value: tt.slug}})
			}

			if err := h.DeletePage(c); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("status mismatch\n  want: %d\n  got:  %d\n  body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}
