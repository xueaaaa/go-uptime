package scheduler

import (
	"context"
	"time"

	"gihub.com/xueaaaa/go-uptime/internal/job"
	"gihub.com/xueaaaa/go-uptime/internal/site/service"
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

func (s *Scheduler) Run(ctx context.Context) error {
	duration := 10 * time.Second /* TODO: config */
	t := time.NewTimer(duration)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:

		}
	}
}

func (s *Scheduler) enqueue(ctx context.Context) error {
	sites, err := s.siteService.GetAll(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, site := range sites {
		if site.NextCheckAt.Before(now) {
			site.LastCheckAt = &now
			site.NextCheckAt = now.Add(site.Interval)
			if err = s.siteService.Update(ctx, site); err != nil {
				return err /* TODO: Log instead of return */
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
