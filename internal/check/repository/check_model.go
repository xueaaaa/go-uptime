package repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CheckModel struct {
	ID         pgtype.UUID
	SiteID     pgtype.UUID
	Status     int
	StatusCode int
	Latency    int
	Error      string
	CheckedAt  time.Time
}
