package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	authmiddleware "github.com/iankencruz/webbuilder/internal/middleware"
	"github.com/iankencruz/webbuilder/internal/session"
)

func (s *Server) registerRoutes() {
	s.e.Use(session.LoadAndSave(s.sessionManager))
	s.e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	auth := s.e.Group("/auth")
	auth.POST("/register", s.authHandler.Register)
	auth.POST("/login", s.authHandler.Login)
	auth.POST("/logout", s.authHandler.Logout)

	s.e.GET("/me", s.authHandler.Me, authmiddleware.RequireAuth(s.sessionManager))
}
