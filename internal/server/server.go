package server

import (
	"context"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/iankencruz/webbuilder/internal/auth"
	"github.com/iankencruz/webbuilder/internal/config"
	"github.com/iankencruz/webbuilder/internal/database/repository"
	"github.com/iankencruz/webbuilder/internal/handler"
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

	queries := repository.New(pool)

	e.Use(middleware.Recover())

	// CORS Strategy Setup
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
	}))

	// Force the Session State engine to wrap around every execution request context cleanly
	e.Use(echo.WrapMiddleware(sessionManager.LoadAndSave))

	authRegistry := initAuthRegistry(cfg)

	authHandler := handler.NewAuthHandler(
		service.NewAuthService(queries),
		sessionManager,
		authRegistry,
	)

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

func initAuthRegistry(cfg *config.Config) *auth.Registry {
	registry := auth.NewRegistry()

	// Register Zitadel provider first (priority)
	if cfg.Zitadel.ClientID != "" && cfg.Zitadel.ClientSecret != "" {
		zitadelProvider := auth.NewOAuth2Provider(
			"zitadel",
			cfg.Zitadel.ClientID,
			cfg.Zitadel.ClientSecret,
			cfg.Zitadel.RedirectURI,
			cfg.Zitadel.IssuerURI,
			cfg.Zitadel.PostLoginURI,
			[]string{"openid", "profile", "email"},
		)
		registry.Register("zitadel", zitadelProvider)
		log.Printf("Registered Zitadel provider with redirect URI: %s", zitadelProvider.RedirectURI())
	}

	if cfg.Google.ClientID != "" && cfg.Google.ClientSecret != "" {
		googleProvider := auth.NewOAuth2Provider(
			"google",
			cfg.Google.ClientID,
			cfg.Google.ClientSecret,
			cfg.Google.RedirectURI,
			"", // Google doesn't need an issuer URI for userinfo
			"", // No post-login redirect for Google
			[]string{"openid", "profile", "email"},
		)
		registry.Register("google", googleProvider)
		log.Printf("Registered Google provider with redirect URI: %s", googleProvider.RedirectURI())
	}

	return registry
}
