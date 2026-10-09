package repository

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xueaaaa/go-uptime/internal/errors"
)

type CheckRepository interface {
	Create(ctx context.Context, check CheckModel) (pgtype.UUID, error)
	Get(ctx context.Context, ID pgtype.UUID) (CheckModel, error)
	GetBySiteID(ctx context.Context, siteID pgtype.UUID) ([]CheckModel, error)
	GetAll(ctx context.Context) ([]CheckModel, error)
	Delete(ctx context.Context, ID pgtype.UUID) error
	DeleteOld(ctx context.Context) error
}

type checkRepository struct {
	db *pgxpool.Pool
}

func NewCheckRepository(db *pgxpool.Pool) CheckRepository {
	return &checkRepository{
		db: db,
	}
}

func (r *checkRepository) Create(ctx context.Context, check CheckModel) (pgtype.UUID, error) {
	sql := `INSERT INTO checks (site_id, status, status_code, latency, error, checked_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id;`

	var id pgtype.UUID
	err := r.db.QueryRow(
		ctx,
		sql,
		check.SiteID,
		check.Status,
		check.StatusCode,
		check.Latency,
		check.Error,
		check.CheckedAt,
	).Scan(&id)

	if err != nil {
		return pgtype.UUID{}, err
	}

	return id, nil
}

func (r *checkRepository) get(ctx context.Context, fieldName string, fieldValue any) ([]CheckModel, error) {
	sql := `SELECT id, site_id, status, status_code, latency, error, checked_at FROM checks`
	var args []any
	if fieldName != "" {
		sql += ` WHERE ` + strings.TrimSpace(fieldName) + ` = $1`
		args = append(args, fieldValue)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checks := make([]CheckModel, 0)
	for rows.Next() {
		check := CheckModel{}
		err = rows.Scan(
			&check.ID,
			&check.SiteID,
			&check.Status,
			&check.StatusCode,
			&check.Latency,
			&check.Error,
			&check.CheckedAt,
		)

		if err != nil {
			return nil, err
		}

		checks = append(checks, check)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return checks, nil
}

func (r *checkRepository) Get(ctx context.Context, ID pgtype.UUID) (CheckModel, error) {
	checks, err := r.get(ctx, "id", ID)
	if err != nil {
		return CheckModel{}, err
	}

	if len(checks) != 1 {
		return CheckModel{}, errors.NotFound
	}

	return checks[0], nil
}

func (r *checkRepository) GetBySiteID(ctx context.Context, siteID pgtype.UUID) ([]CheckModel, error) {
	checks, err := r.get(ctx, "site_id", siteID)
	if err != nil {
		return nil, err
	}

	if len(checks) == 0 {
		return nil, errors.NotFound
	}

	return checks, nil
}

func (r *checkRepository) GetAll(ctx context.Context) ([]CheckModel, error) {
	checks, err := r.get(ctx, "", "")
	if err != nil {
		return nil, err
	}

	if len(checks) == 0 {
		return nil, errors.NoChecks
	}

	return checks, nil
}

func (r *checkRepository) Delete(ctx context.Context, ID pgtype.UUID) error {
	sql := `DELETE FROM checks WHERE id = $1`

	tag, err := r.db.Exec(ctx, sql, ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.NotFound
	}
	return nil
}

func (r *checkRepository) DeleteOld(ctx context.Context) error {
	sql := `DELETE FROM checks WHERE checked_at < now() - interval '30 days'`

	_, err := r.db.Exec(ctx, sql)
	if err != nil {
		return err
	}
	return nil
}
