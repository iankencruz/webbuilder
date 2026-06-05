package service

import (
	"log/slog"

	"github.com/iankencruz/webbuilder/internal/blocks"
	"github.com/iankencruz/webbuilder/internal/database/repository"
)

type Services struct {
	Auth  *AuthService
	Page  *PageService
	Block *BlockService
}

func NewServices(
	logger *slog.Logger,
	queries *repository.Queries,
) *Services {
	return &Services{
		Auth:  NewAuthService(logger, queries),
		Page:  NewPageService(logger, queries),
		Block: NewBlockService(logger, queries, []blocks.BlockType{blocks.RichText}),
	}
}
