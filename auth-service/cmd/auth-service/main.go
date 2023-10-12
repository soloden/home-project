package main

import (
	pkgApp "auth-service/internal/app"
	"context"
	"fmt"
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
	app, err := pkgApp.NewApp(context.TODO())
	if err != nil {
		return fmt.Errorf("initialization app: %s", err)
	}

	if err := app.Run(); err != nil {
		return fmt.Errorf("running app: %s", err)
	}

	// listen to OS signals and gracefully shutdown
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
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
