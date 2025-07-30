package utils

import (
	"time"
)

func ParseTimeISO(timeStr string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func FormatTimeISO(t time.Time) string {
	return t.Format(time.RFC3339)
}
