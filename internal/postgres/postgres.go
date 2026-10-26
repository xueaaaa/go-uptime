package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xueaaaa/go-uptime/internal/config"
)

func NewConn(ctx context.Context, conf config.Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DB,
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	config.MinConns = conf.MinConns
	config.MaxConns = conf.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
