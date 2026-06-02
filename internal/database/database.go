package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

func Setup(ctx context.Context, databaseURL, migrationsDir string) (*pgxpool.Pool, error) {
	// Initialize primary pgx pool
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// ping to ensure connection is valid
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations using goose
	sqlDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to open sql DB: %w", err)
	}
	defer sqlDB.Close()

	if err := RunMigrations(sqlDB, migrationsDir); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return pool, nil
}

func RunMigrations(sqlDB *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(sqlDB, migrationsDir)
}
