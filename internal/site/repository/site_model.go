package repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type SiteModel struct {
	ID               pgtype.UUID
	URL              string
	Status           int
	ConsecutiveFails int
	Interval         time.Duration
	LastCheckAt      *time.Time
	NextCheckAt      time.Time
	CreatedAt        time.Time
}
