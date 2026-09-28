package model

import (
	"time"

	model2 "gihub.com/xueaaaa/go-uptime/internal/site/model"
	"github.com/google/uuid"
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
