package handler

import (
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v5"

	"github.com/iankencruz/webbuilder/internal/auth"
	"github.com/iankencruz/webbuilder/internal/database/repository"
)

type AuthServicer interface {
	FindOrCreateUser(ctx echo.Context, sub, email, firstName, lastName, avatarURL, provider string) (repository.User, error)
	GetByID(ctx echo.Context, id int64) (repository.User, error)
}

type AuthHandler struct {
	logger       *slog.Logger
	authService  AuthServicer
	sessions     *scs.SessionManager
	authRegistry *auth.Registry
}

func NewAuthHandler(log *slog.Logger, authService AuthServicer, sessions *scs.SessionManager, oidcRegistry *auth.Registry) *AuthHandler {
	return &AuthHandler{
		logger:       log,
		authService:  authService,
		sessions:     sessions,
		authRegistry: oidcRegistry,
	}
}

func (h *AuthHandler) OAuthLogin(c *echo.Context) error {
	providerName := c.Param("provider")
	if providerName == "" {
		providerName = c.Param("provider")
	}

	h.logger.Info("auth login request", "provider", providerName)

	provider, ok := h.authRegistry.Get(providerName)
	if !ok {
		h.logger.Error("provider not found", "provider", providerName)
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

	h.logger.Info("session initialized", "state", state)

	url := provider.AuthenticationURL(state)

	h.logger.Info("redirecting to IDP", "url", url)

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) OAuthCallback(c *echo.Context) error {
	ctx := c.Request().Context()

	// retrieve and pop state, nonce, provider from session
	storedState := h.sessions.PopString(ctx, "oauth_state")
	providerName := h.sessions.PopString(ctx, "oauth_provider")
	incomingState := c.QueryParam("state")
	incomingCode := c.QueryParam("code")

	h.logger.Info("callback session", "stored_state", storedState, "provider", providerName)
	h.logger.Info("callback params", "incoming_state", incomingState, "code_len", len(incomingCode))

	if providerName == "" {
		providerName = c.Param("provider")
		h.logger.Debug("provider fallback", "provider", providerName)
	}

	if incomingState != storedState || storedState == "" {
		h.logger.Error("state mismatch", "incoming", incomingState, "stored", storedState)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "invalid anti-forgery token state tracking"})
	}

	provider, ok := h.authRegistry.Get(providerName)
	if !ok {
		h.logger.Error("provider not found", "provider", providerName)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "unknown provider"})
	}

	// exchange code for tokens
	profile, err := provider.ExchangeCode(ctx, c.QueryParam("code"))
	if err != nil {
		h.logger.Error("token exchange failed", "err", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "failed to exchange token"})
	}

	user, err := h.authService.FindOrCreateUser(
		ctx,
		profile.Subject,
		profile.Email,
		profile.GivenName,
		profile.FamilyName,
		profile.Avatar,
		"",
	)
	if err != nil {
		h.logger.Error("find or create user failed", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to find or create user"})
	}

	// create session
	h.sessions.Put(ctx, "user_id", user.ID)
	postLoginURL := provider.PostLoginURL()
	h.logger.Info("callback success", "redirect", postLoginURL)
	return c.Redirect(http.StatusTemporaryRedirect, postLoginURL)
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
	h.logger.Info("AUTH: Me Handler",
		"user_found", err == nil,
		slog.Group("User Information",
			slog.Int64("id", user.ID),
			slog.String("email", user.Email),
			slog.String("first name", user.FirstName.String),
			slog.String("last name", user.LastName.String),
		))

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// Internal Helper Utility Methods
func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
