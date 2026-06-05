package main

import (
	"context"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/iankencruz/webbuilder/internal/config"
	"github.com/iankencruz/webbuilder/internal/database"
	"github.com/iankencruz/webbuilder/internal/server"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	pool, err := database.Setup(ctx, cfg.DatabaseURL, cfg.MigrationsDir)
	if err != nil {
		log.Fatalf("setup database: %v", err)
	}
	defer pool.Close()

	app := server.New(ctx, cfg, pool)

	if err := app.Start(ctx); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
