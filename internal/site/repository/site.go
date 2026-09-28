package repository

import (
	"context"

	"gihub.com/xueaaaa/go-uptime/internal/errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SiteRepository interface {
	Create(ctx context.Context, site SiteModel) (pgtype.UUID, error)
	Get(ctx context.Context, ID pgtype.UUID) (SiteModel, error)
	Update(ctx context.Context, site SiteModel) error
	Delete(ctx context.Context, ID pgtype.UUID) error
}

type siteRepository struct {
	db *pgxpool.Pool
}

func NewSiteRepository(db *pgxpool.Pool) SiteRepository {
	return &siteRepository{
		db: db,
	}
}

func (r *siteRepository) Create(ctx context.Context, site SiteModel) (pgtype.UUID, error) {
	sql := `INSERT INTO sites (url, status, consecutive_fails, interval, last_check_at, next_check_at, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id;`

	var id pgtype.UUID
	err := r.db.QueryRow(
		ctx,
		sql,
		site.URL,
		site.Status,
		site.ConsecutiveFails,
		site.Interval,
		site.LastCheckAt,
		site.NextCheckAt,
		site.CreatedAt,
	).Scan(&id)

	if err != nil {
		return pgtype.UUID{}, err
	}

	return id, nil
}

func (r *siteRepository) Get(ctx context.Context, ID pgtype.UUID) (SiteModel, error) {
	sql := `SELECT id, url, status, consecutive_fails, interval, last_check_at, next_check_at, created_at FROM sites
			WHERE id = $1`

	var site SiteModel
	err := r.db.QueryRow(ctx, sql, ID).Scan(
		&site.ID,
		&site.URL,
		&site.Status,
		&site.ConsecutiveFails,
		&site.Interval,
		&site.LastCheckAt,
		&site.NextCheckAt,
		&site.CreatedAt,
	)
	if err != nil {
		return SiteModel{}, errors.NotFound
	}

	return site, nil
}

func (r *siteRepository) Update(ctx context.Context, site SiteModel) error {
	sql := `UPDATE sites
			SET url =  COALESCE(NULLIF($1::text, ''), url), 
			    status = COALESCE(NULLIF($2::int, 0), status),
			    consecutive_fails = $3,
			    interval = $4,
			    last_check_at = $5,
			    next_check_at = $6
			    WHERE id = $7`

	tag, err := r.db.Exec(
		ctx,
		sql,
		site.URL,
		site.Status,
		site.ConsecutiveFails,
		site.Interval,
		site.LastCheckAt,
		site.NextCheckAt,
		site.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.NotFound
	}

	return nil
}

func (r *siteRepository) Delete(ctx context.Context, ID pgtype.UUID) error {
	sql := `DELETE FROM sites WHERE id = $1`

	tag, err := r.db.Exec(ctx, sql, ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.NotFound
	}
	return nil
}
