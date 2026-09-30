package job

import "github.com/google/uuid"

type Job struct {
	SiteID uuid.UUID
	URL    string
}
