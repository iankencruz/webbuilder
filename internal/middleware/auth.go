package middleware

import (
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
)

func RequireAuth(sessionManager *scs.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := sessionManager.GetString(c.Request().Context(), "user_id")
			if _, err := strconv.ParseInt(userID, 10, 64); err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			return next(c)
		}
	}
}
