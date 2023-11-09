package app

import (
	"context"
	"time"

	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Error(msg string)
	Info(msg string)
}

type Storage interface {
	Create(ctx context.Context, event *storage.Event) error
	Get(ctx context.Context, id int) (storage.Event, error)
	Update(ctx context.Context, event *storage.Event) error
	Delete(ctx context.Context, id int) error
	ListDayEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error)
	ListWeekEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error)
	ListMonthEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context,
	id, userID int,
	title, desc string,
	start, notifTime time.Time,
	dur time.Duration,
) error {
	event := storage.Event{
		ID:               id,
		Title:            title,
		StartDate:        start,
		Duration:         dur,
		Description:      desc,
		UserID:           userID,
		NotificationTime: notifTime,
	}
	err := a.Storage.Create(ctx, &event)
	if err != nil {
		return err
	}

	return nil
}
