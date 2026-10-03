package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/xueaaaa/go-uptime/cmd"
	"github.com/xueaaaa/go-uptime/internal/postgres"
	"github.com/xueaaaa/go-uptime/internal/site/repository"
	"github.com/xueaaaa/go-uptime/internal/site/service"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pool, err := postgres.NewConn(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	siteRepo := repository.NewSiteRepository(pool)
	siteSvc := service.NewSiteService(siteRepo)

	return cmd.NewRootCmd(siteSvc).ExecuteContext(ctx)
}
