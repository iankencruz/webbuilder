package service

import (
	"context"
	"log/slog"

	"github.com/iankencruz/webbuilder/internal/database/repository"
)

type PageService struct {
	queries *repository.Queries
}

func NewPageService(logger *slog.Logger, q *repository.Queries) *PageService {
	return &PageService{queries: q}
}

func (s *PageService) CreatePage(ctx context.Context, arg repository.CreatePageParams) (repository.Page, error) {
	return s.queries.CreatePage(ctx, arg)
}

func (s *PageService) GetPageByID(ctx context.Context, id int64) (repository.Page, error) {
	return s.queries.GetPageByID(ctx, id)
}

func (s *PageService) GetPageBySlug(ctx context.Context, slug string) (repository.Page, error) {
	return s.queries.GetPageBySlug(ctx, slug)
}

func (s *PageService) ListPages(ctx context.Context) ([]repository.Page, error) {
	return s.queries.ListPages(ctx)
}

func (s *PageService) UpdatePage(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error) {
	return s.queries.UpdatePage(ctx, arg)
}

func (s *PageService) DeletePageByID(ctx context.Context, id int64) error {
	return s.queries.DeleteByID(ctx, id)
}

func (s *PageService) DeletePageBySlug(ctx context.Context, slug string) error {
	return s.queries.DeleteBySlug(ctx, slug)
}
