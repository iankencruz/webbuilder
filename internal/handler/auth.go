package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v5"
	"golang.org/x/oauth2"

	"github.com/iankencruz/webbuilder/internal/oidc"
	"github.com/iankencruz/webbuilder/internal/service"
)

type AuthHandler struct {
	authService  *service.AuthService
	sessions     *scs.SessionManager
	oidcRegistry *oidc.Registry
}

func NewAuthHandler(authService *service.AuthService, sessions *scs.SessionManager, oidcRegistry *oidc.Registry) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		sessions:     sessions,
		oidcRegistry: oidcRegistry,
	}
}

func (h *AuthHandler) OAuthLogin(c *echo.Context) error {
	providerName := c.Param("provider")

	provider, ok := h.oidcRegistry.Get(providerName)
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "unknown provider"})
	}

	state, err := generateState()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate state"})
	}

	nonce, err := generateState() // same random generation, different purpose
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate nonce"})
	}

	h.sessions.Put(c.Request().Context(), "oauth_state", state)
	h.sessions.Put(c.Request().Context(), "oauth_nonce", nonce)
	h.sessions.Put(c.Request().Context(), "oauth_provider", providerName)

	url := provider.Config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("nonce", nonce),
	)

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) OAuthCallback(c *echo.Context) error {
	ctx := c.Request().Context()

	// retrieve and pop state, nonce, provider from session
	storedState := h.sessions.PopString(ctx, "oauth_state")
	storedNonce := h.sessions.PopString(ctx, "oauth_nonce")
	providerName := h.sessions.PopString(ctx, "oauth_provider")

	// validate state
	if c.QueryParam("state") != storedState || storedState == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid state"})
	}

	provider, ok := h.oidcRegistry.Get(providerName)
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "unknown provider"})
	}

	// exchange code for tokens
	token, err := provider.Config.Exchange(ctx, c.QueryParam("code"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "failed to exchange token"})
	}

	// extract raw id_token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing id_token"})
	}

	// verify id token
	idToken, err := provider.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid id_token"})
	}

	// extract claims
	var claims struct {
		Sub     string `json:"sub"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
		Nonce   string `json:"nonce"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to parse claims"})
	}

	// validate nonce
	if claims.Nonce != storedNonce || storedNonce == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid nonce"})
	}

	// find or create user
	user, err := h.authService.FindOrCreateUser(ctx, claims.Sub, providerName, claims.Email, claims.Name, claims.Picture)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to find or create user"})
	}

	// create session
	h.sessions.Put(ctx, "user_id", user.ID)

	return c.Redirect(http.StatusTemporaryRedirect, provider.PostLoginURL)
}

func (h *AuthHandler) Logout(c *echo.Context) error {
	h.sessions.Remove(c.Request().Context(), "user_id")
	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) Me(c *echo.Context) error {
	userID := h.sessions.GetInt64(c.Request().Context(), "user_id")
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	user, err := h.authService.GetByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, map[string]any{"id": user.ID, "email": user.Email})
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
