package service

import (
	"context"
	"log/slog"

	"github.com/iankencruz/webbuilder/internal/database/repository"
)

type PageRepository interface {
	CreatePage(ctx context.Context, arg repository.CreatePageParams) (repository.Page, error)
	GetPageByID(ctx context.Context, id int64) (repository.Page, error)
	GetPageBySlug(ctx context.Context, slug string) (repository.Page, error)
	ListPages(ctx context.Context) ([]repository.Page, error)
	UpdatePage(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error)
	DeleteByID(ctx context.Context, id int64) error
	DeleteBySlug(ctx context.Context, slug string) error
}

type PageService struct {
	repo PageRepository
}

func NewPageService(logger *slog.Logger, q *repository.Queries) *PageService {
	return &PageService{repo: q}
}

func (s *PageService) CreatePage(ctx context.Context, arg repository.CreatePageParams) (repository.Page, error) {
	return s.repo.CreatePage(ctx, arg)
}

func (s *PageService) GetPageByID(ctx context.Context, id int64) (repository.Page, error) {
	return s.repo.GetPageByID(ctx, id)
}

func (s *PageService) GetPageBySlug(ctx context.Context, slug string) (repository.Page, error) {
	return s.repo.GetPageBySlug(ctx, slug)
}

func (s *PageService) ListPages(ctx context.Context) ([]repository.Page, error) {
	return s.repo.ListPages(ctx)
}

func (s *PageService) UpdatePage(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error) {
	return s.repo.UpdatePage(ctx, arg)
}

func (s *PageService) DeletePageByID(ctx context.Context, id int64) error {
	return s.repo.DeleteByID(ctx, id)
}

func (s *PageService) DeletePageBySlug(ctx context.Context, slug string) error {
	return s.repo.DeleteBySlug(ctx, slug)
}
