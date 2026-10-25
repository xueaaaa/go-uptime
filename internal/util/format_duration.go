package util

import (
	"strings"
	"time"
)

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)

	days := int(d / (24 * time.Hour))
	d %= 24 * time.Hour
	hours := int(d / time.Hour)
	d %= time.Hour
	minutes := int(d / time.Minute)
	d %= time.Minute
	seconds := int(d / time.Second)

	var parts []string
	if days > 0 {
		parts = append(parts, Pluralize(days, "day", "days"))
	}
	if hours > 0 {
		parts = append(parts, Pluralize(hours, "hour", "hours"))
	}
	if minutes > 0 {
		parts = append(parts, Pluralize(minutes, "minute", "minutes"))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, Pluralize(seconds, "second", "seconds"))
	}

	return strings.Join(parts, " ")
}
