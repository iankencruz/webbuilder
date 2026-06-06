package handler

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

type mockPageService struct {
	listPages     func(ctx context.Context) ([]repository.Page, error)
	createPage    func(ctx context.Context, arg repository.CreatePageParams) (repository.Page, error)
	getPageByID   func(ctx context.Context, id int64) (repository.Page, error)
	getPageBySlug func(ctx context.Context, slug string) (repository.Page, error)
	updatePage    func(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error)
	deleteByID    func(ctx context.Context, id int64) error
	deleteBySlug  func(ctx context.Context, slug string) error
}

func (m *mockPageService) ListPages(ctx context.Context) ([]repository.Page, error) {
	return m.listPages(ctx)
}

func (m *mockPageService) CreatePage(ctx context.Context, arg repository.CreatePageParams) (repository.Page, error) {
	return m.createPage(ctx, arg)
}

func (m *mockPageService) GetPageByID(ctx context.Context, id int64) (repository.Page, error) {
	return m.getPageByID(ctx, id)
}

func (m *mockPageService) GetPageBySlug(ctx context.Context, slug string) (repository.Page, error) {
	return m.getPageBySlug(ctx, slug)
}

func (m *mockPageService) UpdatePage(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error) {
	return m.updatePage(ctx, arg)
}

func (m *mockPageService) DeletePageByID(ctx context.Context, id int64) error {
	return m.deleteByID(ctx, id)
}

func (m *mockPageService) DeletePageBySlug(ctx context.Context, slug string) error {
	return m.deleteBySlug(ctx, slug)
}

// --- Helpers ----

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
		slug       string
		svc        *mockPageService
		wantStatus int
		wantLen    int
	}{
		{
			name: "List Pages",
			svc: &mockPageService{
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
			svc: &mockPageService{
				listPages: func(ctx context.Context) ([]repository.Page, error) {
					return []repository.Page{}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantLen:    0,
		},
		{
			name: "service error returns 500",
			svc: &mockPageService{
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
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewPageHandler(logger, tt.svc)
			c, rec := newTestContext(e, http.MethodGet, "/pages", "")

			if err := handler.ListPages(c); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			if tt.wantStatus == http.StatusOK {
				var pages []repository.Page
				if err := json.Unmarshal(rec.Body.Bytes(), &pages); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if len(pages) != tt.wantLen {
					t.Errorf("expected %d pages, got %d", tt.wantLen, len(pages))
				}
			}
		})
	}
}

func TestCreatePage(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svc        *mockPageService
		wantStatus int
		wantTitle  string
		wantSlug   string
	}{
		{
			name: "valid request creates page",
			body: `{"title":"Contact","slug":"contact","status":"draft"}`,
			svc: &mockPageService{
				createPage: func(ctx context.Context, arg repository.CreatePageParams) (repository.Page, error) {
					return repository.Page{ID: 3, Title: arg.Title, Slug: arg.Slug}, nil
				},
			},
			wantStatus: http.StatusCreated,
			wantTitle:  "Contact",
			wantSlug:   "contact",
		},
		{
			name:       "missing slug returns 400",
			body:       `{"title":"Missing Slug","status":"draft"}`,
			svc:        &mockPageService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing title returns 400",
			body:       `{"slug":"missing-title","status":"draft"}`,
			svc:        &mockPageService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error returns 500",
			body: `{"title":"Error Page","slug":"error-page","status":"draft"}`,
			svc: &mockPageService{
				createPage: func(ctx context.Context, arg repository.CreatePageParams) (repository.Page, error) {
					return repository.Page{}, errors.New("db connection refused")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewPageHandler(logger, tt.svc)

			req := httptest.NewRequest(http.MethodPost, "/api/pages", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if err := handler.CreatePage(c); err != nil {
				t.Fatalf("handler returned unexpected error: %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("status mismatch\n  want: %d\n  got:  %d\n  body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}

			if tt.wantTitle != "" || tt.wantSlug != "" {
				var page repository.Page
				if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
					t.Fatalf("failed to decode response body: %v\n  raw body: %s", err, rec.Body.String())
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
		svc        *mockPageService
		wantStatus int
		wantTitle  string
		wantSlug   string
	}{
		{
			name: "valid slug returns page",
			slug: "home",
			svc: &mockPageService{
				getPageBySlug: func(ctx context.Context, slug string) (repository.Page, error) {
					return repository.Page{ID: 1, Title: "Home", Slug: slug}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantTitle:  "Home",
			wantSlug:   "home",
		},
		{
			name:       "missing slug returns 400",
			slug:       "",
			svc:        &mockPageService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "wrong slug returns 404",
			slug: "nonexistent",
			svc: &mockPageService{
				getPageBySlug: func(ctx context.Context, slug string) (repository.Page, error) {
					return repository.Page{}, pgx.ErrNoRows
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "not found returns 404",
			slug: "nonexistent",
			svc: &mockPageService{
				getPageBySlug: func(ctx context.Context, slug string) (repository.Page, error) {
					return repository.Page{}, pgx.ErrNoRows
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "service error returns 500",
			slug: "home",
			svc: &mockPageService{
				getPageBySlug: func(ctx context.Context, slug string) (repository.Page, error) {
					return repository.Page{}, errors.New("db connection refused")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewPageHandler(logger, tt.svc)

			req := httptest.NewRequest(http.MethodGet, "/api/pages/"+tt.slug, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.slug != "" {
				c.SetPathValues(echo.PathValues{
					{Name: "slug", Value: tt.slug},
				})
			}

			if err := handler.GetPage(c); err != nil {
				t.Fatalf("handler returned unexpected error: %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("status mismatch\n  want: %d\n  got:  %d\n  body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}

			if tt.wantTitle != "" || tt.wantSlug != "" {
				var page repository.Page
				if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
					t.Fatalf("failed to decode response body: %v\n  raw body: %s", err, rec.Body.String())
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

func TestUpdatePage(t *testing.T) {
	tests := []struct {
		name       string
		slug       string
		body       string
		svc        *mockPageService
		wantStatus int
		wantTitle  string
		wantSlug   string
	}{
		{
			name: "valid request updates page",
			slug: "home",
			body: `{"title":"Home Updated","status":"published"}`,
			svc: &mockPageService{
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
			body:       `{"title":"Home Updated","status":"published"}`,
			svc:        &mockPageService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "page not found returns 404",
			slug: "nonexistent",
			body: `{"title":"Home Updated","status":"published"}`,
			svc: &mockPageService{
				updatePage: func(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error) {
					return repository.Page{}, pgx.ErrNoRows
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "service error returns 500",
			slug: "home",
			body: `{"title":"Home Updated","status":"published"}`,
			svc: &mockPageService{
				updatePage: func(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error) {
					return repository.Page{}, errors.New("db connection refused")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewPageHandler(logger, tt.svc)

			req := httptest.NewRequest(http.MethodPut, "/api/pages/"+tt.slug, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.slug != "" {
				c.SetPathValues(echo.PathValues{
					{Name: "slug", Value: tt.slug},
				})
			}

			if err := handler.UpdatePage(c); err != nil {
				t.Fatalf("handler returned unexpected error: %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("status mismatch\n  want: %d\n  got:  %d\n  body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}

			if tt.wantTitle != "" || tt.wantSlug != "" {
				var page repository.Page
				if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
					t.Fatalf("failed to decode response body: %v\n  raw body: %s", err, rec.Body.String())
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
		svc        *mockPageService
		wantStatus int
	}{
		{
			name: "valid slug deletes page",
			slug: "home",
			svc: &mockPageService{
				deleteBySlug: func(ctx context.Context, slug string) error {
					return nil
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "page not found returns 404",
			slug: "nonexistent",
			svc: &mockPageService{
				deleteBySlug: func(ctx context.Context, slug string) error {
					return pgx.ErrNoRows
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "service error returns 500",
			slug: "home",
			svc: &mockPageService{
				deleteBySlug: func(ctx context.Context, slug string) error {
					return errors.New("db connection refused")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := NewPageHandler(logger, tt.svc)

			req := httptest.NewRequest(http.MethodDelete, "/api/pages/"+tt.slug, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.slug != "" {
				c.SetPathValues(echo.PathValues{
					{Name: "slug", Value: tt.slug},
				})
			}

			if err := handler.DeletePage(c); err != nil {
				t.Fatalf("handler returned unexpected error: %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("status mismatch\n  want: %d\n  got:  %d\n  body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}
