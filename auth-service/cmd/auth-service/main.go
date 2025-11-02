package main

import (
	pkgApp "auth-service/internal/app"
	"auth-service/internal/config"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}

	os.Exit(1)
}

func run() error {
	fmt.Println("starting auth service...")
	cfg := config.MustLoad()
	log := setupLogger(cfg.App.ENV)

	app, err := pkgApp.NewApp(log)
	if err != nil {
		return fmt.Errorf("initialization app failed: %s", err)
	}

	if err := app.Run(); err != nil {
		return fmt.Errorf("running app failed: %s", err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("It is ok"))
		if err != nil {
			log.Error("body write error: %s", err)
		}
	})

	log.Info("start listening HTTP server on the 8080 port")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return fmt.Errorf("problem with start HTTP server: %s", err)
	}

	// listen to OS signals and gracefully shutdown
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 3)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := app.Stop(ctx); err != nil {
			fmt.Printf("gracefull shutdown problem: %s", err)
		}
		defer cancel()
		close(stopped)
	}()
	<-stopped
	return nil
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case "production":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	case "development":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case "local":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	return logger
}
