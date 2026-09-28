package model

type Status int

const (
	Unknown     Status = 0
	Unavailable Status = 1
	Available   Status = 2
)
