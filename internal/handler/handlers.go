package handler

import (
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/iankencruz/webbuilder/internal/auth"
	"github.com/iankencruz/webbuilder/internal/service"
)

type Handlers struct {
	Auth  *AuthHandler
	Page  *PageHandler
	Block *BlockHandler
}

func NewHandler(
	logger *slog.Logger,
	services *service.Services,
	sessions *scs.SessionManager,
	authRegistry *auth.Registry,
) *Handlers {
	return &Handlers{
		Auth:  NewAuthHandler(logger, services.Auth, sessions, authRegistry),
		Page:  NewPageHandler(logger, services.Page),
		Block: NewBlockHandler(logger, services.Block),
	}
}
