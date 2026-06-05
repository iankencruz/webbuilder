package server

import (
	"net/http"

	"github.com/labstack/echo/v5"

	authmiddleware "github.com/iankencruz/webbuilder/internal/middleware"
)

func (s *Server) registerRoutes() {
	s.e.Use(authmiddleware.RequestLogger())

	s.e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
	})

	api := s.e.Group("/api")

	auth := api.Group("/auth")
	auth.POST("/logout", s.handlers.Auth.Logout)
	auth.GET("/login/:provider", s.handlers.Auth.OAuthLogin)
	auth.GET("/callback/:provider", s.handlers.Auth.OAuthCallback)

	api.GET("/me", s.handlers.Auth.Me, authmiddleware.RequireAuth(s.sessionManager))
}
