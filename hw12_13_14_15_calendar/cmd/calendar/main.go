package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/config"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/soadmized/go_otus/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.json", "path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logg := logger.New(conf.LogLevel)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	restAPI := internalhttp.NewServer(logg, calendar, conf.Addr())
	grpcAPI := internalgrpc.NewServer(logg, calendar, conf.GRPCAddr())

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := restAPI.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		if err := grpcAPI.Stop(ctx); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	go func() {
		if err := grpcAPI.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	if err := restAPI.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
