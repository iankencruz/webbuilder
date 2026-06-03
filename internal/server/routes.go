package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	authmiddleware "github.com/iankencruz/webbuilder/internal/middleware"
	"github.com/iankencruz/webbuilder/pkg/logger"
)

func (s *Server) registerRoutes() {
	s.e.Use(middleware.RequestLoggerWithConfig(logger.RequestLoggerConfig()))

	s.e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
	})

	api := s.e.Group("/api")

	auth := api.Group("/auth")
	auth.POST("/logout", s.authHandler.Logout)
	auth.GET("/login/:provider", s.authHandler.OAuthLogin)
	auth.GET("/callback/:provider", s.authHandler.OAuthCallback)

	api.GET("/me", s.authHandler.Me, authmiddleware.RequireAuth(s.sessionManager))
}
