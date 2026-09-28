package service

import (
	"context"

	"gihub.com/xueaaaa/go-uptime/internal/site/model"
	"gihub.com/xueaaaa/go-uptime/internal/site/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SiteService interface {
	Create(ctx context.Context, site model.Site) (uuid.UUID, error)
	Get(ctx context.Context, ID uuid.UUID) (model.Site, error)
	Update(ctx context.Context, site model.Site) error
	Delete(ctx context.Context, ID uuid.UUID) error
}

type siteService struct {
	repo repository.SiteRepository
}

func NewCheckService(repo repository.SiteRepository) SiteService {
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

func (s *siteService) Delete(ctx context.Context, ID uuid.UUID) error {
	pgID := pgtype.UUID{
		Bytes: ID,
		Valid: true,
	}
	return s.repo.Delete(ctx, pgID)
}
