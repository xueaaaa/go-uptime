package service

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/xueaaaa/go-uptime/internal/check/model"
	model2 "github.com/xueaaaa/go-uptime/internal/site/model"
	"github.com/xueaaaa/go-uptime/internal/site/repository"
)

type fakeRepo struct {
	site repository.SiteModel
}

func (f *fakeRepo) Get(_ context.Context, _ pgtype.UUID) (repository.SiteModel, error) {
	return f.site, nil
}

func (f *fakeRepo) Update(_ context.Context, site repository.SiteModel) error {
	f.site = site
	return nil
}

func (f *fakeRepo) Create(context.Context, repository.SiteModel) (pgtype.UUID, error) {
	panic("unexpected Create")
}
func (f *fakeRepo) GetByUrl(context.Context, string) (repository.SiteModel, error) {
	panic("unexpected GetByUrl")
}
func (f *fakeRepo) GetAll(context.Context) ([]repository.SiteModel, error) {
	panic("unexpected GetAll")
}
func (f *fakeRepo) Delete(context.Context, pgtype.UUID) error {
	panic("unexpected Delete")
}

func TestTwoFailuresMakeSiteUnavailable(t *testing.T) {
	repo := &fakeRepo{
		site: repository.SiteModel{
			Status:   int(model2.Available),
			Interval: time.Minute,
		},
	}
	svc := NewSiteService(repo)

	failedCheck := model.Check{
		Status: model2.Unavailable,
	}

	if err := svc.UpdateByCheck(context.Background(), failedCheck); err != nil {
		t.Fatal(err)
	}
	if repo.site.Status != int(model2.Available) {
		t.Fatal("site became unavailable first after the failures")
	}

	if err := svc.UpdateByCheck(context.Background(), failedCheck); err != nil {
		t.Fatal(err)
	}
	if repo.site.Status != int(model2.Unavailable) {
		t.Fatal("site did not become unavailable after the second failure")
	}
	if repo.site.ConsecutiveFails != 2 {
		t.Fatalf("expected 2 failures, got %d", repo.site.ConsecutiveFails)
	}
}

func TestSuccessMakesSiteAvailable(t *testing.T) {
	repo := &fakeRepo{
		site: repository.SiteModel{
			Status:           int(model2.Unavailable),
			ConsecutiveFails: 2,
			Interval:         time.Minute,
		},
	}
	svc := NewSiteService(repo)

	successCheck := model.Check{
		Status: model2.Available,
	}

	if err := svc.UpdateByCheck(context.Background(), successCheck); err != nil {
		t.Fatal(err)
	}
	if repo.site.Status != int(model2.Available) {
		t.Fatal("site did not become accessible after successful check")
	}
	if repo.site.ConsecutiveFails != 0 {
		t.Fatal("after a successful check, the number of failed checks was not reset to 0")
	}
}
