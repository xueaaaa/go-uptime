package scheduler

import (
	"context"
	"errors"
	"time"

	"github.com/xueaaaa/go-uptime/internal/config"
	appErrors "github.com/xueaaaa/go-uptime/internal/errors"
	"github.com/xueaaaa/go-uptime/internal/job"
	"github.com/xueaaaa/go-uptime/internal/site/service"
)

type Scheduler struct {
	siteService service.SiteService
	jobs        chan<- job.Job
}

func NewScheduler(
	siteService service.SiteService,
	jobs chan<- job.Job,
) *Scheduler {
	return &Scheduler{
		siteService: siteService,
		jobs:        jobs,
	}
}

func (s *Scheduler) Run(ctx context.Context, conf config.Config) error {
	duration := conf.TickInterval
	t := time.NewTicker(duration)
	defer t.Stop()
	defer close(s.jobs)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			if err := s.enqueue(ctx); err != nil {
				return err
			}
		}
	}
}

func (s *Scheduler) enqueue(ctx context.Context) error {
	sites, err := s.siteService.GetAll(ctx)
	if errors.Is(err, appErrors.NoSites) {
		return nil
	}
	if err != nil {
		return err
	}

	now := time.Now()
	for _, site := range sites {
		if site.NextCheckAt.Before(now) {
			site.LastCheckAt = &now
			site.NextCheckAt = now.Add(site.Interval)
			if err = s.siteService.Update(ctx, site); err != nil {
				return err
			}

			select {
			case s.jobs <- job.Job{
				SiteID: site.ID,
				URL:    site.URL,
			}:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}
