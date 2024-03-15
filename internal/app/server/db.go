package server

import (
	"context"
	"level0/internal/app/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getPool(ctx context.Context, config *config.Config) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, config.DBAddress)
}

func initTables(ctx context.Context, config *config.Config) error {
	query := `
		create table if not exists "order" (
			id serial primary key,
			error_message varchar,
			order_info_json jsonb
		)
	`
	_, err := config.Pool.Exec(ctx, query)
	return err
}
