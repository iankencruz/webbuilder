package server

import (
	"net/http"

	"github.com/labstack/echo/v5"

	authmiddleware "github.com/iankencruz/webbuilder/internal/middleware"
)

func (s *Server) registerRoutes() {
	s.e.Use(authmiddleware.RequestLogger())

	s.e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
	})

	api := s.e.Group("/api")

	auth := api.Group("/auth")
	api.GET("/me", s.handlers.Auth.Me, authmiddleware.RequireAuth(s.sessionManager))
	auth.POST("/logout", s.handlers.Auth.Logout)
	auth.GET("/login/:provider", s.handlers.Auth.OAuthLogin)
	auth.GET("/callback/:provider", s.handlers.Auth.OAuthCallback)

	// Pages
	pages := api.Group("/pages", authmiddleware.RequireAuth(s.sessionManager))
	pages.GET("", s.handlers.Page.ListPages)
	pages.POST("", s.handlers.Page.CreatePage)
	pages.GET("/:id", s.handlers.Page.GetPage)
	pages.PUT("/:id", s.handlers.Page.UpdatePage)
	pages.DELETE("/:id", s.handlers.Page.DeletePage)

	// Page Blocks (junction)
	blocks := pages.Group("/:id/blocks")
	blocks.POST("", s.handlers.Block.AddBlockToPage)
	blocks.PUT("/:block_id", s.handlers.Block.UpdatePageBlock)
	blocks.DELETE("/:block_id", s.handlers.Block.DeletePageBlock)

	// Blocks (type-specific)
	block := api.Group("/blocks", authmiddleware.RequireAuth(s.sessionManager))
	block.POST("/:collection", s.handlers.Block.CreateBlock)
	block.GET("/:collection/:id", s.handlers.Block.GetBlock)
	block.PUT("/:collection/:id", s.handlers.Block.UpdateBlock)
}
