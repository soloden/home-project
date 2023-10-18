package main

import (
	pkgApp "auth-service/internal/app"
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mode := flag.String("mode", "", "that flag declares mode to run service")
	flag.Parse()
	if err := run(*mode); err != nil {
		panic(err)
	}

	os.Exit(1)
}

func run(mode string) error {
	log := zap.Must(zap.NewProduction())
	if mode == "development" {
		log = zap.Must(zap.NewDevelopment())
	}

	app, err := pkgApp.NewApp(context.TODO())
	if err != nil {
		return fmt.Errorf("initialization app: %s", err)
	}

	if err := app.Run(); err != nil {
		return fmt.Errorf("running app: %s", err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("It is ok"))
		if err != nil {
			log.Error("body write error: %s", zap.Error(err))
		}
	})

	fmt.Println("start listening HTTP server on the 8080 port")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return fmt.Errorf("problem with start HTTP server: %s", err)
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
