package handler

import (
	"log/slog"
	"net/http"

	"github.com/iankencruz/webbuilder/internal/database/repository"
	"github.com/iankencruz/webbuilder/internal/service"
	"github.com/labstack/echo/v5"
)

type PageHandler struct {
	logger   *slog.Logger
	services *service.PageService
}

func NewPageHandler(log *slog.Logger, services *service.PageService) *PageHandler {
	return &PageHandler{
		logger:   log,
		services: services,
	}
}

func (h *PageHandler) ListPages(c *echo.Context) error {
	pages, err := h.services.ListPages(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list pages", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list pages"})
	}
	return c.JSON(http.StatusOK, pages)
}

func (h *PageHandler) CreatePage(c *echo.Context) error {
	var params repository.CreatePageParams
	if err := c.Bind(&params); err != nil {
		h.logger.Error("failed to bind request body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if params.Title == "" || params.Slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title and slug are required"})
	}

	page, err := h.services.CreatePage(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("failed to create page", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create page"})
	}
	return c.JSON(http.StatusCreated, page)
}

func (h *PageHandler) GetPage(c *echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "slug is required"})
	}

	page, err := h.services.GetPageBySlug(c.Request().Context(), slug)
	if err != nil {
		h.logger.Error("failed to get page", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get page"})
	}
	return c.JSON(http.StatusOK, page)
}

func (h *PageHandler) UpdatePage(c *echo.Context) error {
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

	page, err := h.services.UpdatePage(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("failed to update page", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update page"})
	}
	return c.JSON(http.StatusOK, page)
}

func (h *PageHandler) DeletePage(c *echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "slug is required"})
	}

	err := h.services.DeletePageBySlug(c.Request().Context(), slug)
	if err != nil {
		h.logger.Error("failed to delete page", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete page"})
	}
	return c.NoContent(http.StatusNoContent)
}
