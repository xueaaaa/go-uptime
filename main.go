package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/pterm/pterm"
	"github.com/xueaaaa/go-uptime/cmd"
	repository2 "github.com/xueaaaa/go-uptime/internal/check/repository"
	service2 "github.com/xueaaaa/go-uptime/internal/check/service"
	"github.com/xueaaaa/go-uptime/internal/config"
	"github.com/xueaaaa/go-uptime/internal/postgres"
	"github.com/xueaaaa/go-uptime/internal/site/repository"
	"github.com/xueaaaa/go-uptime/internal/site/service"
)

func main() {
	if err := run(); err != nil {
		pterm.Error.Print(err)
		os.Exit(1)
	}
}

func run() error {
	conf, err := config.New()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pool, err := postgres.NewConn(ctx, conf)
	if err != nil {
		return err
	}
	defer pool.Close()

	siteRepo := repository.NewSiteRepository(pool)
	siteSvc := service.NewSiteService(siteRepo)
	checkRepo := repository2.NewCheckRepository(pool)
	checkSvc := service2.NewCheckService(checkRepo)

	return cmd.NewRootCmd(siteSvc, checkSvc, conf).ExecuteContext(ctx)
}
