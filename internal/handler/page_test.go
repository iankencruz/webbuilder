package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iankencruz/webbuilder/internal/database/repository"
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

func newTestContext(e *echo.Context, method, target, body string) (echo.Context, *httptest.ResponseRecorder) {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, target, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, target, nil)
	}

	rec := httptest.NewRecorder()
	return *e.Echo().NewContext(req, rec), rec
}

// --- Tests ---

func TestListPages(t *testing.T) {
}
