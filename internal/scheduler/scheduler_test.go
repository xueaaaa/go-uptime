package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	model2 "github.com/xueaaaa/go-uptime/internal/check/model"
	"github.com/xueaaaa/go-uptime/internal/job"
	"github.com/xueaaaa/go-uptime/internal/site/model"
)

type fakeService struct {
	site model.Site
}

func (s *fakeService) GetAll(_ context.Context) ([]model.Site, error) {
	return []model.Site{s.site}, nil
}

func (s *fakeService) Update(_ context.Context, site model.Site) error { s.site = site; return nil }

func (s *fakeService) Create(context.Context, model.Site) (uuid.UUID, error) {
	panic("unexpected Create")
}
func (s *fakeService) Get(context.Context, uuid.UUID) (model.Site, error) { panic("unexpected Get") }
func (s *fakeService) GetByUrl(context.Context, string) (model.Site, error) {
	panic("unexpected GetByUrl")
}
func (s *fakeService) UpdateByCheck(context.Context, model2.Check) error {
	panic("unexpected UpdateByCheck")
}
func (s *fakeService) Delete(context.Context, uuid.UUID) error { panic("unexpected Delete") }

func TestSiteReady(t *testing.T) {
	site := model.Site{
		ID:          uuid.New(),
		URL:         "https://example.com",
		Interval:    10 * time.Second,
		NextCheckAt: time.Now().Add(-time.Second),
	}
	service := &fakeService{
		site: site,
	}
	jobs := make(chan job.Job, 1)

	sched := NewScheduler(service, jobs)
	before := time.Now()
	if err := sched.enqueue(context.Background()); err != nil {
		t.Fatal(err)
	}
	after := time.Now()

	select {
	case got := <-jobs:
		if got.SiteID != site.ID {
			t.Fatalf("expected %v SiteID, got %v", site.ID, got.SiteID)
		}
		if got.URL != site.URL {
			t.Fatalf("expected %q URL, got %q", site.URL, got.URL)
		}
		next := service.site.NextCheckAt
		if next.Before(before.Add(site.Interval)) ||
			next.After(after.Add(site.Interval)) {
			t.Fatalf("invalid NextCheckTime: %v", next)
		}
	default:
		t.Fatalf("expected 1 job, got 0")
	}
}

func TestSiteNotReady(t *testing.T) {
	site := model.Site{
		NextCheckAt: time.Now().Add(time.Minute),
	}
	service := &fakeService{
		site: site,
	}
	jobs := make(chan job.Job, 1)

	sched := NewScheduler(service, jobs)
	if err := sched.enqueue(context.Background()); err != nil {
		t.Fatal(err)
	}

	select {
	case <-jobs:
		t.Fatalf("expected 0 jobs, got 1")
	default:
		if !service.site.NextCheckAt.Equal(site.NextCheckAt) {
			t.Fatal("check time was changed ahead of schedule")
		}
	}
}
