package blocks

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/iankencruz/webbuilder/internal/database/repository"
	"github.com/labstack/echo/v5"
)

type BlockServicer interface {
	Resolve(collection string) (Block, error)
	AddBlockToPage(ctx context.Context, arg repository.AddBlockToPageParams) (repository.PagesBlock, error)
	UpdatePageBlock(ctx context.Context, arg repository.UpdatePageBlockParams) (repository.PagesBlock, error)
	DeletePageBlock(ctx context.Context, junctionID int64) error
	ReorderPageBlocks(ctx context.Context, arg repository.ReorderPageBlocksParams) error
}

type BlockHandler struct {
	logger   *slog.Logger
	services BlockServicer
}

func NewBlockHandler(log *slog.Logger, services BlockServicer) *BlockHandler {
	return &BlockHandler{
		logger:   log,
		services: services,
	}
}

func (h *BlockHandler) CreateBlock(c *echo.Context) error {
	collection := c.Param("collection")

	block, err := h.services.Resolve(collection)
	if err != nil {
		h.logger.Error("failed to resolve block collection", "collection", collection, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to resolve block collection"})
	}

	if err := c.Bind(block); err != nil {
		h.logger.Error("failed to bind request body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	id, err := block.Create(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to create block", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create block"})
	}

	return c.JSON(http.StatusCreated, map[string]any{"id": id})
}

func (h *BlockHandler) GetBlock(c *echo.Context) error {
	collection := c.Param("type")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("invalid block ID", "id", c.Param("id"), "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid block ID"})
	}

	block, err := h.services.Resolve(collection)
	if err != nil {
		h.logger.Error("failed to resolve block collection", "collection", collection, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to resolve block collection"})
	}

	result, err := block.Get(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to get block", "id", id, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get block"})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *BlockHandler) UpdateBlock(c *echo.Context) error {
	collection := c.Param("type")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("invalid block ID", "id", c.Param("id"), "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid block ID"})
	}

	block, err := h.services.Resolve(collection)
	if err != nil {
		h.logger.Error("failed to resolve block collection", "collection", collection, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to resolve block collection"})
	}

	if err := c.Bind(block); err != nil {
		h.logger.Error("failed to bind request body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	result, err := block.Update(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to update block", "id", id, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update block"})
	}

	return c.JSON(http.StatusOK, result)
}

// --- Junction Handlers ---
func (h *BlockHandler) AddBlockToPage(c *echo.Context) error {
	pageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("invalid page ID", "id", c.Param("id"), "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid page ID"})
	}

	var params repository.AddBlockToPageParams
	if err := c.Bind(&params); err != nil {
		h.logger.Error("failed to bind request body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	params.PageID = pageID

	junction, err := h.services.AddBlockToPage(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("failed to add block to page", "page_id", pageID, "block_id", params.BlockID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to add block to page"})
	}

	return c.JSON(http.StatusOK, junction)
}

func (h *BlockHandler) UpdatePageBlock(c *echo.Context) error {
	junctionID, err := strconv.ParseInt(c.Param("junctionID"), 10, 64)
	if err != nil {
		h.logger.Error("invalid junction ID", "id", c.Param("junctionID"), "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid junction ID"})
	}

	var params repository.UpdatePageBlockParams
	if err := c.Bind(&params); err != nil {
		h.logger.Error("failed to bind request body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	params.ID = junctionID

	junction, err := h.services.UpdatePageBlock(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("failed to update page block", "junction_id", junctionID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update page block"})
	}

	return c.JSON(http.StatusOK, junction)
}

func (h *BlockHandler) DeletePageBlock(c *echo.Context) error {
	junctionID, err := strconv.ParseInt(c.Param("junctionID"), 10, 64)
	if err != nil {
		h.logger.Error("invalid junction ID", "id", c.Param("junctionID"), "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid junction ID"})
	}

	if err := h.services.DeletePageBlock(c.Request().Context(), junctionID); err != nil {
		h.logger.Error("failed to delete page block", "junction_id", junctionID, "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete page block"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *BlockHandler) ReorderPageBlocks(c *echo.Context) error {
	pageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("invalid page ID", "id", c.Param("id"), "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid page ID"})
	}
	var items []repository.ReorderPageBlocksParams
	if err := c.Bind(&items); err != nil {
		h.logger.Error("failed to bind request body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	for _, item := range items {
		if err := h.services.ReorderPageBlocks(c.Request().Context(), item); err != nil {
			h.logger.Error("failed to reorder page blocks", "page_id", pageID, "block_id", item.ID, "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to reorder page blocks"})
		}
	}

	return c.NoContent(http.StatusNoContent)
}
