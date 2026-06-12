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
	api.GET("/me", s.handlers.auth.Me, authmiddleware.RequireAuth(s.sessionManager))
	auth.POST("/logout", s.handlers.auth.Logout)
	auth.GET("/login/:provider", s.handlers.auth.OAuthLogin)
	auth.GET("/callback/:provider", s.handlers.auth.OAuthCallback)

	// Pages
	pages := api.Group("/pages", authmiddleware.RequireAuth(s.sessionManager))
	pages.GET("", s.handlers.page.ListPages)
	pages.POST("", s.handlers.page.CreatePage)
	pages.GET("/:slug", s.handlers.page.GetPage)
	pages.PUT("/:slug", s.handlers.page.UpdatePage)
	pages.DELETE("/:slug", s.handlers.page.DeletePage)

	// Page Blocks (junction)
	blocks := pages.Group("/:id/blocks")
	blocks.POST("", s.handlers.block.AddBlockToPage)
	blocks.PUT("/:block_id", s.handlers.block.UpdatePageBlock)
	blocks.DELETE("/:block_id", s.handlers.block.DeletePageBlock)

	// Blocks (type-specific)
	block := api.Group("/blocks", authmiddleware.RequireAuth(s.sessionManager))
	block.POST("/:collection", s.handlers.block.CreateBlock)
	block.GET("/:collection/:id", s.handlers.block.GetBlock)
	block.PUT("/:collection/:id", s.handlers.block.UpdateBlock)
}
