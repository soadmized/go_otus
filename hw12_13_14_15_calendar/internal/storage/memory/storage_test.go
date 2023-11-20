package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) { //nolint:funlen
	t.Parallel()

	t.Run("positive crud", func(t *testing.T) {
		ctx := context.Background()

		start := time.Now().Add(time.Hour * 24)
		nTime := time.Now().Add(time.Hour * 24 * 2)

		e1 := storage.Event{
			ID:               1,
			Title:            "first",
			StartDate:        start,
			Duration:         time.Hour * 24 * 3,
			Description:      "first event",
			UserID:           1,
			NotificationTime: nTime,
		}

		e2 := storage.Event{
			ID:               1,
			Title:            "first updated",
			StartDate:        start,
			Duration:         time.Hour * 24 * 3,
			Description:      "updated first event",
			UserID:           1,
			NotificationTime: nTime,
		}

		s := New()

		err := s.Create(ctx, &e1)
		require.NoError(t, err)

		event, err := s.Get(ctx, e1.ID)
		require.NoError(t, err)
		require.Equal(t, e1, event)

		err = s.Update(ctx, &e2)
		require.NoError(t, err)
		event, err = s.Get(ctx, e1.ID)
		require.NoError(t, err)
		require.Equal(t, e2, event)

		err = s.Delete(ctx, e1.ID)
		require.NoError(t, err)
		_, err = s.Get(ctx, e1.ID)
		require.Error(t, err)
	})

	t.Run("create event already exists", func(t *testing.T) {
		ctx := context.Background()

		start := time.Now().Add(time.Hour * 24)
		nTime := time.Now().Add(time.Hour * 24 * 2)

		e1 := storage.Event{
			ID:               1,
			Title:            "first",
			StartDate:        start,
			Duration:         time.Hour * 24 * 3,
			Description:      "first event",
			UserID:           1,
			NotificationTime: nTime,
		}

		e2 := storage.Event{
			ID:               1,
			Title:            "second",
			StartDate:        start,
			Duration:         time.Hour * 24 * 3,
			Description:      "second event",
			UserID:           2,
			NotificationTime: nTime,
		}

		s := New()

		err := s.Create(ctx, &e1)
		require.NoError(t, err)

		err = s.Create(ctx, &e2)
		require.Error(t, err)
	})

	t.Run("update not existent", func(t *testing.T) {
		ctx := context.Background()

		start := time.Now().Add(time.Hour * 24)
		nTime := time.Now().Add(time.Hour * 24 * 2)

		e1 := storage.Event{
			ID:               1,
			Title:            "first",
			StartDate:        start,
			Duration:         time.Hour * 24 * 3,
			Description:      "first event",
			UserID:           1,
			NotificationTime: nTime,
		}

		e2 := storage.Event{
			ID:               2,
			Title:            "second",
			StartDate:        start,
			Duration:         time.Hour * 24 * 3,
			Description:      "second event",
			UserID:           2,
			NotificationTime: nTime,
		}

		s := New()

		err := s.Create(ctx, &e1)
		require.NoError(t, err)

		err = s.Update(ctx, &e2)
		require.Error(t, err)
	})

	t.Run("delete not found", func(t *testing.T) {
		ctx := context.Background()

		start := time.Now().Add(time.Hour * 24)
		nTime := time.Now().Add(time.Hour * 24 * 2)

		e1 := storage.Event{
			ID:               1,
			Title:            "first",
			StartDate:        start,
			Duration:         time.Hour * 24 * 3,
			Description:      "first event",
			UserID:           1,
			NotificationTime: nTime,
		}

		s := New()

		err := s.Create(ctx, &e1)
		require.NoError(t, err)

		err = s.Delete(ctx, 2)
		require.Error(t, err)
	})

	t.Run("list day events", func(t *testing.T) {
		ctx := context.Background()

		e1 := storage.Event{
			ID:               1,
			Title:            "first",
			StartDate:        time.Now().Add(time.Hour * 20),
			Duration:         time.Hour * 24 * 3,
			Description:      "first event",
			UserID:           1,
			NotificationTime: time.Now().Add(time.Hour * 24 * 2),
		}

		e2 := storage.Event{
			ID:               2,
			Title:            "second",
			StartDate:        time.Now().Add(time.Hour * 24 * 2),
			Duration:         time.Hour * 24 * 3,
			Description:      "second event",
			UserID:           2,
			NotificationTime: time.Now().Add(time.Hour * 24 * 3),
		}

		s := New()

		err := s.Create(ctx, &e1)
		require.NoError(t, err)
		err = s.Create(ctx, &e2)
		require.NoError(t, err)

		event, err := s.ListDayEvents(ctx, time.Now().Add(time.Hour*20))
		require.NoError(t, err)
		require.Len(t, event, 1)
		require.Equal(t, 1, event[0].ID)
	})

	t.Run("list week events", func(t *testing.T) {
		ctx := context.Background()

		e1 := storage.Event{
			ID:               1,
			Title:            "first",
			StartDate:        time.Now().Add(time.Hour * 24 * 3),
			Duration:         time.Hour * 24 * 3,
			Description:      "first event",
			UserID:           1,
			NotificationTime: time.Now().Add(time.Hour * 24 * 5),
		}

		e2 := storage.Event{
			ID:               2,
			Title:            "second",
			StartDate:        time.Now().Add(time.Hour * 24 * 20),
			Duration:         time.Hour * 24 * 3,
			Description:      "second event",
			UserID:           2,
			NotificationTime: time.Now().Add(time.Hour * 24 * 22),
		}

		s := New()

		err := s.Create(ctx, &e1)
		require.NoError(t, err)
		err = s.Create(ctx, &e2)
		require.NoError(t, err)

		event, err := s.ListWeekEvents(ctx, time.Now().Add(time.Hour*24*2))
		require.NoError(t, err)
		require.Len(t, event, 1)
		require.Equal(t, 1, event[0].ID)
	})

	t.Run("list month events", func(t *testing.T) {
		ctx := context.Background()

		e1 := storage.Event{
			ID:               1,
			Title:            "first",
			StartDate:        time.Now().Add(time.Hour * 24 * 13),
			Duration:         time.Hour * 24 * 3,
			Description:      "first event",
			UserID:           1,
			NotificationTime: time.Now().Add(time.Hour * 24 * 14),
		}

		e2 := storage.Event{
			ID:               2,
			Title:            "second",
			StartDate:        time.Now().Add(time.Hour * 24 * 42),
			Duration:         time.Hour * 24 * 3,
			Description:      "second event",
			UserID:           2,
			NotificationTime: time.Now().Add(time.Hour * 24 * 44),
		}

		s := New()

		err := s.Create(ctx, &e1)
		require.NoError(t, err)
		err = s.Create(ctx, &e2)
		require.NoError(t, err)

		event, err := s.ListMonthEvents(ctx, time.Now().Add(time.Hour*24))
		require.NoError(t, err)
		require.Len(t, event, 1)
		require.Equal(t, 1, event[0].ID)
	})
}
