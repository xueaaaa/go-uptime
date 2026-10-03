package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/xueaaaa/go-uptime/internal/check/model"
	"github.com/xueaaaa/go-uptime/internal/check/repository"
	model2 "github.com/xueaaaa/go-uptime/internal/site/model"
)

type CheckService interface {
	Create(ctx context.Context, check model.Check) (uuid.UUID, error)
	Get(ctx context.Context, ID uuid.UUID) (model.Check, error)
	GetBySiteID(ctx context.Context, siteID uuid.UUID) ([]model.Check, error)
	Delete(ctx context.Context, ID uuid.UUID) error
}

type checkService struct {
	repo repository.CheckRepository
}

func NewCheckService(repo repository.CheckRepository) CheckService {
	return &checkService{
		repo: repo,
	}
}

func (s *checkService) Create(ctx context.Context, check model.Check) (uuid.UUID, error) {
	checkModel := repository.CheckModel{
		SiteID: pgtype.UUID{
			Bytes: check.SiteID,
			Valid: true,
		},
		Status:     int(check.Status),
		StatusCode: check.StatusCode,
		Latency:    check.Latency,
		Error:      check.Error,
		CheckedAt:  check.CheckedAt,
	}

	id, err := s.repo.Create(ctx, checkModel)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.UUID(id.Bytes), nil
}

func (s *checkService) Get(ctx context.Context, ID uuid.UUID) (model.Check, error) {
	pgID := pgtype.UUID{
		Bytes: ID,
		Valid: true,
	}

	checkModel, err := s.repo.Get(ctx, pgID)
	if err != nil {
		return model.Check{}, err
	}

	return model.Check{
		ID:         uuid.UUID(checkModel.ID.Bytes),
		SiteID:     uuid.UUID(checkModel.SiteID.Bytes),
		Status:     model2.Status(checkModel.Status),
		StatusCode: checkModel.StatusCode,
		Latency:    checkModel.Latency,
		Error:      checkModel.Error,
		CheckedAt:  checkModel.CheckedAt,
	}, nil
}

func (s *checkService) GetBySiteID(ctx context.Context, siteID uuid.UUID) ([]model.Check, error) {
	pgID := pgtype.UUID{
		Bytes: siteID,
		Valid: true,
	}

	checkModels, err := s.repo.GetBySiteID(ctx, pgID)
	if err != nil {
		return nil, err
	}

	checks := make([]model.Check, len(checkModels))
	for i, v := range checkModels {
		checks[i] = model.Check{
			ID:         uuid.UUID(v.ID.Bytes),
			SiteID:     uuid.UUID(v.SiteID.Bytes),
			Status:     model2.Status(v.Status),
			StatusCode: v.StatusCode,
			Latency:    v.Latency,
			Error:      v.Error,
			CheckedAt:  v.CheckedAt,
		}
	}

	return checks, nil
}

func (s *checkService) Delete(ctx context.Context, ID uuid.UUID) error {
	pgID := pgtype.UUID{
		Bytes: ID,
		Valid: true,
	}
	return s.repo.Delete(ctx, pgID)
}
