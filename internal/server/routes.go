package server

import (
	"net/http"

	"github.com/labstack/echo/v5"

	authmiddleware "github.com/iankencruz/webbuilder/internal/middleware"
)

func (s *Server) registerRoutes() {
	s.e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
	})

	api := s.e.Group("/api")

	auth := api.Group("/auth")
	auth.POST("/logout", s.authHandler.Logout)
	auth.GET("/:provider/login", s.authHandler.OAuthLogin)
	auth.GET("/:provider/callback", s.authHandler.OAuthCallback)

	api.GET("/me", s.authHandler.Me, authmiddleware.RequireAuth(s.sessionManager))
}
