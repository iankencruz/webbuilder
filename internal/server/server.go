package server

import (
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/iankencruz/webbuilder/internal/config"
	"github.com/iankencruz/webbuilder/internal/handler"
)

type Server struct {
	e              *echo.Echo
	cfg            config.Config
	authHandler    *handler.AuthHandler
	sessionManager *scs.SessionManager
}

func New(cfg config.Config, authHandler *handler.AuthHandler, sessionManager *scs.SessionManager) *Server {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	s := &Server{
		e:              e,
		cfg:            cfg,
		authHandler:    authHandler,
		sessionManager: sessionManager,
	}
	s.registerRoutes()

	return s
}

func (s *Server) Start() error {
	return s.e.Start(s.cfg.Addr)
}
