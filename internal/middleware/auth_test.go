package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
)

func TestRequireAuth_UnauthorizedWithoutSession(t *testing.T) {
	e := echo.New()
	sessionManager := scs.New()

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	ctx, err := sessionManager.Load(req.Context(), "")
	if err != nil {
		t.Fatalf("load session: %v", err)
	}
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := RequireAuth(sessionManager)(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	if err := h(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestRequireAuth_AllowsAuthenticatedSession(t *testing.T) {
	e := echo.New()
	sessionManager := scs.New()

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	ctx, err := sessionManager.Load(req.Context(), "")
	if err != nil {
		t.Fatalf("load session: %v", err)
	}
	sessionManager.Put(ctx, "user_id", "1")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := RequireAuth(sessionManager)(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	if err := h(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
