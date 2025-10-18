package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill)
	defer stop()

	r := chi.NewMux()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello world!")
	})

	// This thing counts every second
	// go func() {
	// 	count := 1
	// 	for {
	// 		slog.Debug(fmt.Sprint("Count ", count))
	// 	}
	// }()

	serverr := make(chan error, 1)
	go func() {
		slog.Debug("Listening to server", "port", port)
		serverr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverr:
		{
			slog.Error("something went wrong", "reason", err.Error())
		}
	case <-ctx.Done():
		{
			if err := ctx.Err(); err != nil {
				slog.Debug("Something went wrong", "error", err.Error())
			}
		}
	}
}
