package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neildavies/swaledale/internal/config"
	"github.com/neildavies/swaledale/internal/seed"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := seed.Run(ctx, pool); err != nil {
		slog.Error("seed database", "error", err)
		os.Exit(1)
	}
	slog.Info("seeded swaledale poc data")
}
