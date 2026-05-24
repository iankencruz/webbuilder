package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v5"
)

func TestRequireAuth(t *testing.T) {
	sessionManager := scs.New()

	tests := []struct {
		name     string
		session  func(ctx context.Context) context.Context
		expected int
	}{
		{
			name: "Unauthorized without session",
			session: func(ctx context.Context) context.Context {
				// Return context as-is with an empty session
				return ctx
			},
			expected: http.StatusUnauthorized,
		},
		{
			name: "Allows authenticated session",
			session: func(ctx context.Context) context.Context {
				// Inject the required user_id into the session context
				sessionManager.Put(ctx, "user_id", int64(1))
				return ctx
			},
			expected: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/me", nil)

			// Load a blank token session into the context
			ctx, err := sessionManager.Load(req.Context(), "")
			if err != nil {
				t.Fatalf("load session: %v", err)
			}

			// Apply case-specific session modifications
			ctx = tt.session(ctx)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Simple dummy handler that returns 200 OK if middleware lets it pass
			h := RequireAuth(sessionManager)(func(c *echo.Context) error {
				return c.NoContent(http.StatusOK)
			})

			if err := h(c); err != nil {
				t.Fatalf("handler error: %v", err)
			}

			if rec.Code != tt.expected {
				t.Errorf("expected status %d, got %d", tt.expected, rec.Code)
			}
		})
	}
}
