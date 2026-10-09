package model

import (
	"time"

	"github.com/google/uuid"
	model2 "github.com/xueaaaa/go-uptime/internal/site/model"
)

type Check struct {
	ID         uuid.UUID
	SiteID     uuid.UUID
	Status     model2.Status
	StatusCode int
	Latency    int
	Error      string
	CheckedAt  time.Time
}
