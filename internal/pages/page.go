package pages

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/iankencruz/webbuilder/internal/database/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v5"
)

type Repository interface {
	CreatePage(ctx context.Context, arg repository.CreatePageParams) (repository.Page, error)
	GetPageByID(ctx context.Context, id int64) (repository.Page, error)
	GetPageBySlug(ctx context.Context, slug string) (repository.Page, error)
	ListPages(ctx context.Context) ([]repository.Page, error)
	UpdatePage(ctx context.Context, arg repository.UpdatePageParams) (repository.Page, error)
	DeleteByID(ctx context.Context, id int64) error
	DeleteBySlug(ctx context.Context, slug string) error
}

type Handler struct {
	logger *slog.Logger
	repo   Repository
}

func NewHandler(log *slog.Logger, repo Repository) *Handler {
	return &Handler{
		logger: log,
		repo:   repo,
	}
}

func (h *Handler) ListPages(c *echo.Context) error {
	pages, err := h.repo.ListPages(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list pages", "error", err)
		fmt.Printf("Error listing pages: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list pages"})
	}
	return c.JSON(http.StatusOK, pages)
}

func (h *Handler) CreatePage(c *echo.Context) error {
	var params repository.CreatePageParams

	if err := c.Bind(&params); err != nil {
		h.logger.Error("failed to bind request body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if params.Title == "" || params.Slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title and slug are required"})
	}

	page, err := h.repo.CreatePage(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("failed to create page", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create page", "details": err.Error()})
	}
	return c.JSON(http.StatusCreated, page)
}

func (h *Handler) GetPage(c *echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "slug is required"})
	}

	page, err := h.repo.GetPageBySlug(c.Request().Context(), slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "page not found"})
		}

		h.logger.Error("failed to get page", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get page"})
	}

	return c.JSON(http.StatusOK, page)
}

func (h *Handler) UpdatePage(c *echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "slug is required"})
	}

	var params repository.UpdatePageParams
	if err := c.Bind(&params); err != nil {
		h.logger.Error("failed to bind request body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	params.Slug = slug

	page, err := h.repo.UpdatePage(c.Request().Context(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "page not found"})
		}

		h.logger.Error("failed to update page", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update page"})
	}
	return c.JSON(http.StatusOK, page)
}

func (h *Handler) DeletePage(c *echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "slug is required"})
	}

	err := h.repo.DeleteBySlug(c.Request().Context(), slug)
	if err != nil {

		h.logger.Error("failed to delete page", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete page"})
	}
	return c.NoContent(http.StatusNoContent)
}
