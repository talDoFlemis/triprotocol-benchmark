package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	healthgo "github.com/hellofresh/health-go/v5"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()
	retcode := 1

	defer func() {
		os.Exit(retcode)
	}()

	slog.InfoContext(ctx, "Loading config")
	settings, err := LoadConfig[Settings]("API", BaseSettings)
	if err != nil {
		slog.ErrorContext(ctx, "failed to load config", slog.Any("err", err))
		return
	}

	health, err := healthgo.New(
		healthgo.WithComponent(healthgo.Component{
			Name:    settings.App.Name,
			Version: settings.App.Version,
		}),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create health checker", slog.Any("err", err))
		return
	}

	router := NewMainHandler(&settings.HTTP, &settings.App, health)

	errChan := make(chan error, 1)
	go func() {
		errChan <- router.Start()
	}()

	select {
	case err = <-errChan:
		slog.ErrorContext(ctx, "error when running server", slog.Any("err", err))
		return
	case <-ctx.Done():
		// Wait for first Signal arrives
	}

	err = router.Shutdown(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to shutdown gracefully the server", slog.Any("err", err))
		return
	}

	slog.InfoContext(ctx, "App stopped gracefully")
	retcode = 0

}
