package server

import (
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/iankencruz/webbuilder/internal/config"
	"github.com/iankencruz/webbuilder/internal/handler"
	"github.com/iankencruz/webbuilder/internal/oidc"
)

type Server struct {
	e              *echo.Echo
	cfg            *config.Config
	authHandler    *handler.AuthHandler
	sessionManager *scs.SessionManager
	oidcRegistry   *oidc.Registry
}

func New(cfg *config.Config, authHandler *handler.AuthHandler, sessionManager *scs.SessionManager, oidcRegistry *oidc.Registry) *Server {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	s := &Server{
		e:              e,
		cfg:            cfg,
		authHandler:    authHandler,
		sessionManager: sessionManager,
		oidcRegistry:   oidcRegistry,
	}
	s.registerRoutes()

	return s
}

func (s *Server) Start() error {
	return s.e.Start(s.cfg.Addr)
}
