package storage

import "time"

type Event struct {
	ID               int
	Title            string
	StartDate        time.Time
	Duration         time.Duration
	Description      string
	UserID           int
	NotificationTime time.Time
}
