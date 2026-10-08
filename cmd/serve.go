package cmd

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
	"github.com/xueaaaa/go-uptime/internal/check/model"
	service2 "github.com/xueaaaa/go-uptime/internal/check/service"
	"github.com/xueaaaa/go-uptime/internal/job"
	"github.com/xueaaaa/go-uptime/internal/processor"
	"github.com/xueaaaa/go-uptime/internal/scheduler"
	"github.com/xueaaaa/go-uptime/internal/site/service"
	"github.com/xueaaaa/go-uptime/internal/worker"
	"golang.org/x/sync/errgroup"
)

func NewServeCmd(siteSvc service.SiteService, checkSvc service2.CheckService) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start background sites monitoring",
		RunE: func(c *cobra.Command, _ []string) error {
			jobs := make(chan job.Job, 100)
			results := make(chan model.Check, 100)

			sched := scheduler.NewScheduler(siteSvc, jobs)
			workPool := worker.NewPool(jobs, results)
			proc := processor.NewProcessor(results, siteSvc, checkSvc)

			g, ctx := errgroup.WithContext(c.Context())
			g.Go(func() error { return sched.Run(ctx) })
			g.Go(func() error { return workPool.Run(ctx) })
			g.Go(func() error { return proc.Run(ctx) })

			err := g.Wait()
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		},
	}
}
