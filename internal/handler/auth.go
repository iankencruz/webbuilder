package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"

	"github.com/iankencruz/webbuilder/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	sessions    *scs.SessionManager
}

func NewAuthHandler(authService *service.AuthService, sessions *scs.SessionManager) *AuthHandler {
	return &AuthHandler{authService: authService, sessions: sessions}
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req authRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email and password are required"})
	}

	user, err := h.authService.Register(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.JSON(http.StatusConflict, map[string]string{"error": "email already registered"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to register"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"id": user.ID, "email": user.Email})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req authRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	user, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	h.sessions.Put(c.Request().Context(), "user_id", strconv.FormatInt(user.ID, 10))
	return c.JSON(http.StatusOK, map[string]interface{}{"id": user.ID, "email": user.Email})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	h.sessions.Remove(c.Request().Context(), "user_id")
	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) Me(c echo.Context) error {
	userIDValue := h.sessions.GetString(c.Request().Context(), "user_id")
	userID, err := strconv.ParseInt(userIDValue, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	user, err := h.authService.GetByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"id": user.ID, "email": user.Email})
}
