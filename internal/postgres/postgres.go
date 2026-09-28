package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConn(ctx context.Context) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("DB"),
	) /* TODO: config */

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	config.MinConns = 5  /* TODO: config */
	config.MaxConns = 25 /* TODO: config */

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
