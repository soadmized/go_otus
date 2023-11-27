package internalgrpc

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/server/grpc/protogen"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestCalendarServer(t *testing.T) { //nolint:funlen
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()

		db := getStorage()
		server := grpcServerStub(ctx, db)

		defer server.Stop(ctx)

		id := 1

		client := grpcClientStub(ctx)
		resp, err := client.Create(ctx, &protogen.CreateRequest{
			Event: &protogen.Event{
				Id:               int64(id),
				Title:            "some event",
				StartDate:        "2023-12-23T18:25:43.511Z",
				Duration:         1000000000000,
				Description:      "sooome event",
				UserId:           1,
				NotificationTime: "2023-12-23T20:25:43.511Z",
			},
		})
		require.NoError(t, err)

		event, err := db.Get(ctx, int(resp.GetId()))
		require.NoError(t, err)
		require.Equal(t, event.ID, int(resp.GetId()))
	})
	t.Run("get", func(t *testing.T) {
		ctx := context.Background()

		db := getStorage()
		server := grpcServerStub(ctx, db)

		defer server.Stop(ctx)

		id := 2

		err := db.Create(ctx, &storage.Event{
			ID:               id,
			Title:            "1",
			StartDate:        time.Now(),
			Duration:         111110,
			Description:      "123",
			UserID:           42,
			NotificationTime: time.Now().Add(time.Hour),
		})
		require.NoError(t, err)

		client := grpcClientStub(ctx)
		resp, err := client.Get(ctx, &protogen.GetRequest{
			Id: int64(2),
		})
		require.NoError(t, err)
		require.Equal(t, int64(2), resp.GetEvent().GetId())
	})
	t.Run("update", func(t *testing.T) {
		ctx := context.Background()

		db := getStorage()
		server := grpcServerStub(ctx, db)

		defer server.Stop(ctx)

		id := 3

		start := time.Now()
		notifTime := start.Add(time.Hour)

		err := db.Create(ctx, &storage.Event{
			ID:               id,
			Title:            "some event",
			StartDate:        start,
			Duration:         111110,
			Description:      "123",
			UserID:           42,
			NotificationTime: notifTime,
		})
		require.NoError(t, err)

		client := grpcClientStub(ctx)
		resp, err := client.Update(ctx, &protogen.UpdateRequest{
			Event: &protogen.Event{
				Id:               3,
				Title:            "updated event",
				StartDate:        start.Format(time.RFC3339),
				Duration:         1000000000000,
				Description:      "sooome event",
				UserId:           42,
				NotificationTime: notifTime.Format(time.RFC3339),
			},
		})
		require.NoError(t, err)

		event, err := db.Get(ctx, int(resp.GetId()))
		require.NoError(t, err)
		require.Equal(t, "updated event", event.Title)
	})
	t.Run("delete", func(t *testing.T) {
		ctx := context.Background()

		db := getStorage()
		server := grpcServerStub(ctx, db)

		defer server.Stop(ctx)

		id := 4

		start := time.Now()
		notifTime := start.Add(time.Hour)

		err := db.Create(ctx, &storage.Event{
			ID:               id,
			Title:            "some event",
			StartDate:        start,
			Duration:         111110,
			Description:      "123",
			UserID:           42,
			NotificationTime: notifTime,
		})
		require.NoError(t, err)

		client := grpcClientStub(ctx)
		resp, err := client.Delete(ctx, &protogen.DeleteRequest{
			Id: int64(4),
		})
		require.NoError(t, err)

		_, err = db.Get(ctx, int(resp.GetId()))
		require.Error(t, err)
	})
	t.Run("list day events", func(t *testing.T) {
		ctx := context.Background()

		db := getStorage()
		server := grpcServerStub(ctx, db)

		defer server.Stop(ctx)

		id := 5
		start := time.Now()
		notifTime := start.Add(time.Hour)

		err := db.Create(ctx, &storage.Event{
			ID:               id,
			Title:            "some event",
			StartDate:        start,
			Duration:         111110,
			Description:      "123",
			UserID:           42,
			NotificationTime: notifTime,
		})
		require.NoError(t, err)

		startTime := time.Now().Add(time.Minute)

		client := grpcClientStub(ctx)
		resp, err := client.ListDayEvents(ctx, &protogen.ListRequest{
			StartDate: startTime.Format(time.RFC3339),
		})
		require.NoError(t, err)
		events := resp.GetEvents()
		require.Len(t, events, 1)
		require.Equal(t, int64(5), events[0].Id)
	})
	t.Run("list week events", func(t *testing.T) {
		ctx := context.Background()

		db := getStorage()
		server := grpcServerStub(ctx, db)

		defer server.Stop(ctx)

		id := 6

		start := time.Now().Add(time.Hour * 24 * 7)
		notifTime := start.Add(time.Hour)

		err := db.Create(ctx, &storage.Event{
			ID:               id,
			Title:            "some event",
			StartDate:        start,
			Duration:         111110,
			Description:      "123",
			UserID:           42,
			NotificationTime: notifTime,
		})
		require.NoError(t, err)

		startTime := time.Now().Add(time.Hour * 24 * 7)

		client := grpcClientStub(ctx)
		resp, err := client.ListWeekEvents(ctx, &protogen.ListRequest{
			StartDate: startTime.Format(time.RFC3339),
		})
		require.NoError(t, err)
		events := resp.GetEvents()
		require.Len(t, events, 1)
		require.Equal(t, int64(6), events[0].Id)
	})
	t.Run("list month events", func(t *testing.T) {
		ctx := context.Background()

		db := getStorage()
		server := grpcServerStub(ctx, db)

		defer server.Stop(ctx)

		id := 7

		start := time.Now().Add(time.Hour * 24 * 31)
		notifTime := start.Add(time.Hour)

		err := db.Create(ctx, &storage.Event{
			ID:               id,
			Title:            "some event",
			StartDate:        start,
			Duration:         111110,
			Description:      "123",
			UserID:           42,
			NotificationTime: notifTime,
		})
		require.NoError(t, err)

		startTime := time.Now().Add(time.Hour * 24 * 31)

		client := grpcClientStub(ctx)
		resp, err := client.ListMonthEvents(ctx, &protogen.ListRequest{
			StartDate: startTime.Format(time.RFC3339),
		})
		require.NoError(t, err)
		events := resp.GetEvents()
		require.Len(t, events, 1)
		require.Equal(t, int64(7), events[0].Id)
	})
}

func grpcServerStub(ctx context.Context, storage *memorystorage.Storage) *Server {
	logg := logger.New("debug")

	application := app.New(logg, storage)
	server := NewServer(logg, application, "localhost:7778")

	go func() {
		if err := server.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	return server
}

func grpcClientStub(ctx context.Context) protogen.CalendarClient {
	conn, err := grpc.DialContext(
		ctx,
		"localhost:7778",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	return protogen.NewCalendarClient(conn)
}

func getStorage() *memorystorage.Storage {
	return memorystorage.New()
}
