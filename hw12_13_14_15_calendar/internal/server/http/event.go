package internalhttp

import (
	"time"

	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/app"
)

type Event struct {
	ID               int           `json:"id,omitempty"`
	Title            string        `json:"title,omitempty"`
	StartDate        time.Time     `json:"startDate,omitempty"`
	Duration         time.Duration `json:"duration,omitempty"`
	Description      string        `json:"description,omitempty"`
	UserID           int           `json:"userID,omitempty"` //nolint:tagliatelle
	NotificationTime time.Time     `json:"notificationTime,omitempty"`
}

func EventToModel(e Event) app.Event {
	return app.Event(e)
}
