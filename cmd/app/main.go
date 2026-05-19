package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/iankencruz/webbuilder/internal/config"
	"github.com/iankencruz/webbuilder/internal/db"
	"github.com/iankencruz/webbuilder/internal/handler"
	"github.com/iankencruz/webbuilder/internal/repository"
	"github.com/iankencruz/webbuilder/internal/server"
	"github.com/iankencruz/webbuilder/internal/service"
	"github.com/iankencruz/webbuilder/internal/session"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("create db pool: %v", err)
	}
	defer pool.Close()

	sqlDB, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open sql db: %v", err)
	}
	defer sqlDB.Close()

	if err := db.RunMigrations(sqlDB, "db/migrations"); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	sessionManager := session.NewManager(pool, cfg.SessionLifetime, cfg.SessionSecure, cfg.SessionCookie)
	queries := repository.New(pool)
	authService := service.NewAuthService(queries)
	authHandler := handler.NewAuthHandler(authService, sessionManager)
	app := server.New(cfg, authHandler, sessionManager)

	if err := app.Start(); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
