package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	model2 "github.com/xueaaaa/go-uptime/internal/check/model"
	"github.com/xueaaaa/go-uptime/internal/site/model"
	"github.com/xueaaaa/go-uptime/internal/site/repository"
)

type SiteService interface {
	Create(ctx context.Context, site model.Site) (uuid.UUID, error)
	Get(ctx context.Context, ID uuid.UUID) (model.Site, error)
	GetByUrl(ctx context.Context, url string) (model.Site, error)
	GetAll(ctx context.Context) ([]model.Site, error)
	Update(ctx context.Context, site model.Site) error
	UpdateByCheck(ctx context.Context, check model2.Check) error
	Delete(ctx context.Context, ID uuid.UUID) error
}

type siteService struct {
	repo repository.SiteRepository
}

func NewSiteService(repo repository.SiteRepository) SiteService {
	return &siteService{
		repo: repo,
	}
}

func (s *siteService) Create(ctx context.Context, site model.Site) (uuid.UUID, error) {
	siteModel := repository.SiteModel{
		URL:              site.URL,
		Status:           int(site.Status),
		ConsecutiveFails: site.ConsecutiveFails,
		Interval:         site.Interval,
		LastCheckAt:      site.LastCheckAt,
		NextCheckAt:      site.NextCheckAt,
		CreatedAt:        site.CreatedAt,
	}

	id, err := s.repo.Create(ctx, siteModel)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.UUID(id.Bytes), nil
}

func (s *siteService) Get(ctx context.Context, ID uuid.UUID) (model.Site, error) {
	pgID := pgtype.UUID{
		Bytes: ID,
		Valid: true,
	}

	siteModel, err := s.repo.Get(ctx, pgID)
	if err != nil {
		return model.Site{}, err
	}

	return model.Site{
		ID:               uuid.UUID(siteModel.ID.Bytes),
		Status:           model.Status(siteModel.Status),
		ConsecutiveFails: siteModel.ConsecutiveFails,
		Interval:         siteModel.Interval,
		LastCheckAt:      siteModel.LastCheckAt,
		NextCheckAt:      siteModel.NextCheckAt,
		CreatedAt:        siteModel.CreatedAt,
	}, nil
}

func (s *siteService) GetByUrl(ctx context.Context, url string) (model.Site, error) {
	siteModel, err := s.repo.GetByUrl(ctx, url)
	if err != nil {
		return model.Site{}, err
	}

	return model.Site{
		ID:               uuid.UUID(siteModel.ID.Bytes),
		Status:           model.Status(siteModel.Status),
		ConsecutiveFails: siteModel.ConsecutiveFails,
		Interval:         siteModel.Interval,
		LastCheckAt:      siteModel.LastCheckAt,
		NextCheckAt:      siteModel.NextCheckAt,
		CreatedAt:        siteModel.CreatedAt,
	}, nil
}

func (s *siteService) GetAll(ctx context.Context) ([]model.Site, error) {
	siteModels, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	sites := make([]model.Site, len(siteModels))
	for i, v := range siteModels {
		site := model.Site{
			ID:               uuid.UUID(v.ID.Bytes),
			URL:              v.URL,
			Status:           model.Status(v.Status),
			ConsecutiveFails: v.ConsecutiveFails,
			Interval:         v.Interval,
			LastCheckAt:      v.LastCheckAt,
			NextCheckAt:      v.NextCheckAt,
			CreatedAt:        v.CreatedAt,
		}

		sites[i] = site
	}

	return sites, nil
}

func (s *siteService) Update(ctx context.Context, site model.Site) error {
	siteModel := repository.SiteModel{
		ID: pgtype.UUID{
			Bytes: site.ID,
			Valid: true,
		},
		URL:              site.URL,
		Status:           int(site.Status),
		ConsecutiveFails: site.ConsecutiveFails,
		Interval:         site.Interval,
		LastCheckAt:      site.LastCheckAt,
		NextCheckAt:      site.NextCheckAt,
		CreatedAt:        site.CreatedAt,
	}

	return s.repo.Update(ctx, siteModel)
}

func (s *siteService) UpdateByCheck(ctx context.Context, check model2.Check) error {
	site, err := s.Get(ctx, check.SiteID)
	if err != nil {
		return err
	}

	if check.Status == model.Unavailable {
		site.ConsecutiveFails++
	} else {
		site.ConsecutiveFails = 0
	}

	if site.ConsecutiveFails >= 2 {
		site.Status = model.Unavailable
	} else if check.Status == model.Available {
		site.Status = model.Available
	}

	now := time.Now()
	site.LastCheckAt = &now

	return s.Update(ctx, site)
}

func (s *siteService) Delete(ctx context.Context, ID uuid.UUID) error {
	pgID := pgtype.UUID{
		Bytes: ID,
		Valid: true,
	}
	return s.repo.Delete(ctx, pgID)
}
