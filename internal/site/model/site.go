package model

import (
	"time"

	"github.com/google/uuid"
)

type Site struct {
	ID               uuid.UUID
	URL              string
	Status           Status
	ConsecutiveFails int
	Interval         time.Duration
	LastCheckAt      *time.Time
	NextCheckAt      time.Time
	CreatedAt        time.Time
}
