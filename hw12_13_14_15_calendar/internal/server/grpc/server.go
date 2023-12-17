package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/server/grpc/protogen"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
)

type App interface {
	Create(ctx context.Context, event app.Event) error
	Get(ctx context.Context, id int) (storage.Event, error)
	Update(ctx context.Context, event app.Event) error
	Delete(ctx context.Context, id int) error
	ListDayEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error)
	ListWeekEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error)
	ListMonthEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error)
}

type Logger interface {
	Error(msg string)
	Info(msg string)
}

type Server struct {
	Server *grpc.Server
	App    App
	Logger Logger
	Addr   string
}

type CalendarServer struct {
	App    App
	Logger Logger
	protogen.UnimplementedCalendarServer
}

func NewServer(logger Logger, app App, addr string) *Server {
	return &Server{
		App:    app,
		Logger: logger,
		Addr:   addr,
	}
}

func (s *Server) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		logMiddleware := loggingMiddleware(s.Logger)
		s.Server = grpc.NewServer(logMiddleware)

		protogen.RegisterCalendarServer(s.Server, &CalendarServer{
			App:    s.App,
			Logger: s.Logger,
		})

		s.Logger.Info(fmt.Sprintf("grpc server start %s", s.Addr))

		listener, err := net.Listen("tcp", s.Addr)
		if err != nil {
			return errors.Wrap(err, "I CAN'T HEAR")
		}

		return s.Server.Serve(listener)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		s.Server.GracefulStop()

		return nil
	}
}

func eventToModel(event *protogen.Event) (app.Event, error) {
	start, err := time.Parse(time.RFC3339, event.GetStartDate())
	if err != nil {
		return app.Event{}, err
	}

	notificationTime, err := time.Parse(time.RFC3339, event.GetNotificationTime())
	if err != nil {
		return app.Event{}, err
	}

	id := int(event.GetId())
	title := event.GetTitle()
	dur := time.Duration(event.GetDuration())
	desc := event.GetDescription()
	userID := int(event.GetUserId())

	return app.Event{
		ID:               id,
		Title:            title,
		StartDate:        start,
		Duration:         dur,
		Description:      desc,
		UserID:           userID,
		NotificationTime: notificationTime,
	}, nil
}

func modelToEvent(event storage.Event) *protogen.Event {
	return &protogen.Event{
		Id:               int64(event.ID),
		Title:            event.Title,
		StartDate:        event.StartDate.String(),
		Duration:         int64(event.Duration),
		Description:      event.Description,
		UserId:           int64(event.UserID),
		NotificationTime: event.NotificationTime.String(),
	}
}
