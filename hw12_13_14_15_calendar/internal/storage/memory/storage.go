package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	storage map[int]storage.Event
	mu      sync.RWMutex
}

func New() *Storage {
	s := make(map[int]storage.Event)

	return &Storage{storage: s}
}

func (s *Storage) Create(_ context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.storage[event.ID]
	if ok {
		msg := fmt.Sprintf("create: event with id=%d already exists", event.ID)

		return errors.New(msg)
	}

	s.storage[event.ID] = *event

	return nil
}

func (s *Storage) Get(_ context.Context, id int) (storage.Event, error) {
	event, ok := s.storage[id]
	if !ok {
		return storage.Event{}, errors.New("event not found")
	}

	return event, nil
}

func (s *Storage) Update(_ context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.storage[event.ID]
	if !ok {
		return errors.New("update: event not found")
	}

	s.storage[event.ID] = *event

	return nil
}

func (s *Storage) Delete(_ context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.storage[id]
	if !ok {
		return errors.New("delete: event not found")
	}

	delete(s.storage, id)

	return nil
}

func (s *Storage) ListDayEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	endDate := startDate.Add(time.Hour * 24)
	events, err := s.listEvents(ctx, endDate)
	if err != nil {
		return nil, errors.Wrap(err, "list day events")
	}

	return events, nil
}

func (s *Storage) ListWeekEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	endDate := startDate.Add(time.Hour * 24 * 7)
	events, err := s.listEvents(ctx, endDate)
	if err != nil {
		return nil, errors.Wrap(err, "list week events")
	}

	return events, nil
}

func (s *Storage) ListMonthEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	endDate := startDate.Add(time.Hour * 24 * 31)
	events, err := s.listEvents(ctx, endDate)
	if err != nil {
		return nil, errors.Wrap(err, "list month events")
	}

	return events, nil
}

func (s *Storage) listEvents(_ context.Context, endDate time.Time) ([]storage.Event, error) { //nolint:unparam
	events := make([]storage.Event, 0)

	for _, v := range s.storage {
		if endDate.After(v.StartDate) {
			events = append(events, v)
		}
	}

	return events, nil
}
