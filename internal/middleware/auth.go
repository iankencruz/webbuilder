package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v5"
)

func RequireAuth(sessionManager *scs.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userID := sessionManager.GetInt64(c.Request().Context(), "user_id")
			if userID == 0 {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			return next(c)
		}
	}
}
