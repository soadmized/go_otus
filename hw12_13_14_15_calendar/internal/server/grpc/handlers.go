package internalgrpc

import (
	"context"
	"time"

	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/server/grpc/protogen"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/storage"
)

func (s *CalendarServer) Create(ctx context.Context, req *protogen.CreateRequest) (*protogen.CreateResponse, error) {
	event, err := eventToModel(req.GetEvent())
	if err != nil {
		return nil, err
	}

	err = s.App.Create(ctx, event)
	if err != nil {
		return nil, err
	}

	return &protogen.CreateResponse{Id: int64(event.ID)}, nil
}

func (s *CalendarServer) Get(ctx context.Context, req *protogen.GetRequest) (*protogen.GetResponse, error) {
	id := int(req.GetId())

	e, err := s.App.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &protogen.GetResponse{Event: modelToEvent(e)}, nil
}

func (s *CalendarServer) Update(ctx context.Context, req *protogen.UpdateRequest) (*protogen.UpdateResponse, error) {
	event, err := eventToModel(req.GetEvent())
	if err != nil {
		return nil, err
	}

	err = s.App.Update(ctx, event)
	if err != nil {
		return nil, err
	}

	return &protogen.UpdateResponse{Id: int64(event.ID)}, nil
}

func (s *CalendarServer) Delete(ctx context.Context, req *protogen.DeleteRequest) (*protogen.DeleteResponse, error) {
	id := int(req.GetId())

	err := s.App.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return &protogen.DeleteResponse{Id: int64(id)}, nil
}

func (s *CalendarServer) ListDayEvents(ctx context.Context, req *protogen.ListRequest) (*protogen.ListResponse, error) {
	resp, err := s.listEvents(ctx, req, s.App.ListDayEvents)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *CalendarServer) ListWeekEvents(ctx context.Context, req *protogen.ListRequest) (*protogen.ListResponse, error) { //nolint:lll
	resp, err := s.listEvents(ctx, req, s.App.ListWeekEvents)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *CalendarServer) ListMonthEvents(ctx context.Context, req *protogen.ListRequest) (*protogen.ListResponse, error) { //nolint:lll
	resp, err := s.listEvents(ctx, req, s.App.ListMonthEvents)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *CalendarServer) listEvents(ctx context.Context, req *protogen.ListRequest, fn func(context.Context, time.Time) ([]storage.Event, error)) (*protogen.ListResponse, error) { //nolint:lll
	start, err := time.Parse(time.RFC3339, req.GetStartDate())
	if err != nil {
		return nil, err
	}

	models, err := fn(ctx, start)
	if err != nil {
		return nil, err
	}

	events := make([]*protogen.Event, 0, len(models))

	for _, m := range models {
		e := modelToEvent(m)
		events = append(events, e)
	}

	return &protogen.ListResponse{Events: events}, nil
}
