package sqlstorage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct { // TODO
	app.Storage
	dsn  string
	conn *pgx.Conn
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Conn(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, s.dsn)
	if err != nil {
		return errors.Wrap(err, "fail to connect db")
	}

	s.conn = conn

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.conn == nil {
		return nil
	}

	err := s.conn.Close(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to close db connection")
	}

	return nil
}

func (s *Storage) Create(ctx context.Context, event *storage.Event) error {
	query := `INSERT INTO event(
                  title, 
                  start_date, 
                  duration, 
                  description,
                  user_id,
                  notification_time
                  ) VALUES ($1, $2, $3, $4) RETURNING id`

	s.conn.QueryRow(
		ctx,
		query,
		event.Title,
		event.StartDate,
		event.Duration,
		event.Description,
		event.UserID,
		event.NotificationTime,
	)

	return nil
}

func (s *Storage) Get(ctx context.Context, id int) (storage.Event, error) {
	e := storage.Event{}
	query := `SELECT id, title, start_date, duration, description,user_id, notification_time FROM e WHERE id=$1`

	row := s.conn.QueryRow(ctx, query, id)

	err := row.Scan(&e.ID, &e.Title, &e.StartDate, &e.Duration, &e.Description, &e.UserID, &e.NotificationTime)
	if err != nil {
		return e, errors.Wrap(err, "get event")
	}

	return e, nil
}

func (s *Storage) Update(ctx context.Context, event *storage.Event) error {
	query := `UPDATE event SET title=$1, start_date=$2, end_date=$3, description=$4, user_id=$5, notification_time=$6 WHERE id=$7` //nolint:lll

	_, err := s.conn.Exec(
		ctx,
		query,
		event.Title,
		event.StartDate,
		event.Duration,
		event.Description,
		event.UserID,
		event.NotificationTime,
		event.ID,
	)
	if err != nil {
		return errors.Wrap(err, "update event")
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM event WHERE id=$1`

	_, err := s.conn.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "delete event")
	}

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
		return nil, errors.Wrap(err, "list day events")
	}

	return events, nil
}

func (s *Storage) listEvents(ctx context.Context, endDate time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)
	query := `SELECT id, title, start_date, duration, description,user_id, notification_time FROM event WHERE start_date >= $1` //nolint:lll

	rows, err := s.conn.Query(ctx, query, endDate)
	if err != nil {
		return nil, errors.Wrap(err, "list events")
	}

	defer rows.Close()

	for rows.Next() {
		var e storage.Event

		err := rows.Scan(&e.ID, &e.Title, &e.StartDate, &e.Duration, &e.Description, &e.UserID, &e.NotificationTime)
		if err != nil {
			return nil, errors.Wrap(err, "list day events")
		}

		events = append(events, e)
	}

	return events, nil
}
