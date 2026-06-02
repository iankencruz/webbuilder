package server

import (
	"context"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/iankencruz/webbuilder/internal/config"
	"github.com/iankencruz/webbuilder/internal/database/repository"
	"github.com/iankencruz/webbuilder/internal/handler"
	"github.com/iankencruz/webbuilder/internal/oidc"
	"github.com/iankencruz/webbuilder/internal/service"
	"github.com/iankencruz/webbuilder/internal/session"
)

type Server struct {
	e              *echo.Echo
	cfg            *config.Config
	authHandler    *handler.AuthHandler
	sessionManager *scs.SessionManager
}

func New(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) *Server {
	e := echo.New()

	sessionManager := session.NewManager(pool, cfg.SessionLifetime, cfg.SessionSecure, cfg.SessionCookie)

	oidcRegistry, err := oidc.NewRegistry(ctx, cfg)
	if err != nil {
		log.Fatalf("build oidc registry: %v", err)
	}

	queries := repository.New(pool)

	authHandler := handler.NewAuthHandler(
		service.NewAuthService(queries),
		sessionManager,
		oidcRegistry,
	)

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
	}))

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
