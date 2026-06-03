package handler

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v5"

	"github.com/iankencruz/webbuilder/internal/auth"
	"github.com/iankencruz/webbuilder/internal/service"
)

type AuthHandler struct {
	authService  *service.AuthService
	sessions     *scs.SessionManager
	authRegistry *auth.Registry
}

func NewAuthHandler(authService *service.AuthService, sessions *scs.SessionManager, oidcRegistry *auth.Registry) *AuthHandler {
	return &AuthHandler{
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

	log.Printf("[AUTH] Init login request for provider: '%s'", providerName)

	provider, ok := h.authRegistry.Get(providerName)
	if !ok {
		log.Printf("[AUTH ERROR] Provider lookup failed inside registry for key: '%s'", providerName)
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

	log.Printf("[AUTH] Session data initialized. State code generated: '%s'", state)

	url := provider.AuthenticationURL(state)

	log.Printf("[AUTH] Redirecting client window out to target IDP URL: %s", url)

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) OAuthCallback(c *echo.Context) error {
	ctx := c.Request().Context()

	log.Printf("[CALLBACK] Endpoint reached! Request URI path: %s", c.Request().RequestURI)

	// retrieve and pop state, nonce, provider from session
	storedState := h.sessions.PopString(ctx, "oauth_state")
	providerName := h.sessions.PopString(ctx, "oauth_provider")
	incomingState := c.QueryParam("state")
	incomingCode := c.QueryParam("code")

	log.Printf("[CALLBACK] Session inspection -> Stored State: '%s', Tracked Provider: '%s'", storedState, providerName)
	log.Printf("[CALLBACK] Query parameters -> Incoming State: '%s', Code length: %d", incomingState, len(incomingCode))

	if providerName == "" {
		providerName = c.Param("provider")
		log.Printf("[CALLBACK WARNING] Session provider tracking was blank, falling back to path segment matching: '%s'", providerName)
	}

	if incomingState != storedState || storedState == "" {
		log.Printf("[CALLBACK FATAL] Security state validation failed! Incoming: '%s', Stored: '%s'", incomingState, storedState)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "invalid anti-forgery token state tracking"})
	}

	provider, ok := h.authRegistry.Get(providerName)
	if !ok {
		log.Printf("[CALLBACK FATAL] Failed registry validation for provider: '%s'", providerName)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "unknown provider"})
	}

	// exchange code for tokens
	profile, err := provider.ExchangeCode(ctx, c.QueryParam("code"))
	if err != nil {
		log.Printf("[CALLBACK FATAL] Core identity token handshake exchange protocol failed: %v", err)
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
		log.Printf("[CALLBACK ERROR] FindOrCreateUser database error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to find or create user"})
	}

	// create session
	h.sessions.Put(ctx, "user_id", user.ID)
	postLoginURL := provider.PostLoginURL()
	log.Printf("[CALLBACK SUCCESS] Redirecting browser to absolute target: '%s'", postLoginURL)
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
	log.Printf("fname=%v lname=%v", user.FirstName, user.LastName)
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
