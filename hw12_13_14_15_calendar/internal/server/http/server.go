package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	app     Application
	logger  Logger
	server  *http.Server
	address string
}

type Logger interface {
	Error(msg string)
	Info(msg string)
}

type Application interface {
	Create(ctx context.Context, event app.Event) error
	Get(ctx context.Context, id int) (storage.Event, error)
	Update(ctx context.Context, event app.Event) error
	Delete(ctx context.Context, id int) error
	ListDayEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error)
	ListWeekEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error)
	ListMonthEvents(ctx context.Context, startDate time.Time) ([]storage.Event, error)
}

func NewServer(logger Logger, app Application, address string) *Server {
	return &Server{
		app:     app,
		logger:  logger,
		address: address,
	}
}

func (s *Server) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		mux := http.NewServeMux()
		handler := Handler{app: s.app, logger: s.logger}

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("hello-world"))
			if err != nil {
				s.logger.Error(err.Error())
			} else {
				w.WriteHeader(200)
			}
		})
		mux.HandleFunc("/create", handler.create)
		mux.HandleFunc("/get", handler.get)
		mux.HandleFunc("/update", handler.update)
		mux.HandleFunc("/delete", handler.delete)
		mux.HandleFunc("/day_list", handler.listDayEvents)
		mux.HandleFunc("/week_list", handler.listWeekEvents)
		mux.HandleFunc("/month_list", handler.listMonthEvents)

		s.server = &http.Server{
			Addr:         s.address,
			Handler:      loggingMiddleware(mux, s.logger),
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}

		s.logger.Info(fmt.Sprintf("listening %s", s.address))

		return s.server.ListenAndServe()
	}
}

func (s *Server) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		return s.server.Close()
	}
}
