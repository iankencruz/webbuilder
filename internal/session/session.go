package session

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func NewManager(pool *pgxpool.Pool, lifetime time.Duration, secure bool) *scs.SessionManager {
	manager := scs.New()
	manager.Store = pgxstore.New(pool)
	manager.Lifetime = lifetime
	manager.Cookie.HttpOnly = true
	manager.Cookie.Persist = true
	manager.Cookie.Secure = secure
	manager.Cookie.SameSite = http.SameSiteLaxMode
	manager.Cookie.Path = "/"
	manager.Cookie.Name = "webbuilder_session"
	return manager
}

func LoadAndSave(manager *scs.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var handlerErr error
			handler := manager.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)
				c.SetResponse(echo.NewResponse(w, c.Echo()))
				handlerErr = next(c)
			}))
			handler.ServeHTTP(c.Response(), c.Request())
			return handlerErr
		}
	}
}
