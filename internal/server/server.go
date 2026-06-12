package server

import (
	"context"
	"log/slog"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/iankencruz/webbuilder/internal/auth"
	"github.com/iankencruz/webbuilder/internal/blocks"
	"github.com/iankencruz/webbuilder/internal/config"
	"github.com/iankencruz/webbuilder/internal/database/repository"
	"github.com/iankencruz/webbuilder/internal/pages"
	"github.com/iankencruz/webbuilder/internal/session"
)

type services struct {
	auth  *auth.AuthService
	page  *pages.PageService
	block *blocks.BlockService
}

type handlers struct {
	auth  *auth.AuthHandler
	page  *pages.PageHandler
	block *blocks.BlockHandler
}

type Server struct {
	e              *echo.Echo
	cfg            *config.Config
	sessionManager *scs.SessionManager
	handlers       *handlers
	services       *services
}

func New(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) *Server {
	e := echo.New()

	var logHandler slog.Handler
	if cfg.Environment == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, nil)
	}
	appLogger := slog.New(logHandler)
	slog.SetDefault(appLogger)

	e.Logger = appLogger

	sessionManager := session.NewManager(pool, cfg.SessionLifetime, cfg.SessionSecure, cfg.SessionCookie)
	queries := repository.New(pool)

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
	}))

	e.Use(echo.WrapMiddleware(sessionManager.LoadAndSave))

	authRegistry := initAuthRegistry(cfg)

	svcs := &services{
		auth:  auth.NewAuthService(appLogger, queries),
		page:  pages.NewPageService(appLogger, queries),
		block: blocks.NewBlockService(appLogger, queries, []blocks.BlockType{blocks.RichText}),
	}

	hdlrs := &handlers{
		auth:  auth.NewAuthHandler(appLogger, svcs.auth, sessionManager, authRegistry),
		page:  pages.NewPageHandler(appLogger, svcs.page),
		block: blocks.NewBlockHandler(appLogger, svcs.block),
	}

	s := &Server{
		e:              e,
		cfg:            cfg,
		sessionManager: sessionManager,
		services:       svcs,
		handlers:       hdlrs,
	}
	s.registerRoutes()

	return s
}

func (s *Server) Start(ctx context.Context) error {
	cfg := echo.StartConfig{
		Address:    s.cfg.Addr,
		HideBanner: true,
		HidePort:   false,
	}

	return cfg.Start(ctx, s.e)
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
	}

	return registry
}
