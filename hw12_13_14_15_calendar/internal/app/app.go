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

type Event struct {
	ID               int
	Title            string
	StartDate        time.Time
	Duration         time.Duration
	Description      string
	UserID           int
	NotificationTime time.Time
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

func (a *App) Create(ctx context.Context, event Event) error {
	e := storage.Event(event)
	err := a.Storage.Create(ctx, &e)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Get(ctx context.Context, id int) (storage.Event, error) {
	event, err := a.Storage.Get(ctx, id)
	if err != nil {
		return storage.Event{}, err
	}

	return event, nil
}

func (a *App) Update(ctx context.Context, event Event) error {
	e := storage.Event(event)
	err := a.Storage.Update(ctx, &e)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Delete(ctx context.Context, id int) error {
	err := a.Storage.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) ListDayEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	events, err := a.Storage.ListDayEvents(ctx, startDate)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (a *App) ListWeekEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	events, err := a.Storage.ListWeekEvents(ctx, startDate)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (a *App) ListMonthEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	events, err := a.Storage.ListMonthEvents(ctx, startDate)
	if err != nil {
		return nil, err
	}

	return events, nil
}
