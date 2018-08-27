package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gyuho/hello-world/version"

	"go.uber.org/zap"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "32001"
	}

	lg, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	lg.Info("starting server",
		zap.String("port", port),
		zap.String("git-commit", version.GitCommit),
		zap.String("release-version", version.ReleaseVersion),
		zap.String("build-time", version.BuildTime),
	)

	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	rootCtx, rootCancel := context.WithCancel(context.Background())
	mux := http.NewServeMux()
	mux.Handle("/hello-world", &contextAdapter{
		lg:      lg,
		ctx:     rootCtx,
		handler: contextHandlerFunc(helloWorldHandler),
	})
	mux.Handle("/hello-world-readiness", &contextAdapter{
		lg:      lg,
		ctx:     rootCtx,
		handler: contextHandlerFunc(readinessHandler),
	})
	mux.Handle("/hello-world-liveness", &contextAdapter{
		lg:      lg,
		ctx:     rootCtx,
		handler: contextHandlerFunc(livenessHandler),
	})
	mux.Handle("/hello-world-status", &contextAdapter{
		lg:      lg,
		ctx:     rootCtx,
		handler: contextHandlerFunc(statusHandler),
	})
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	errc := make(chan error)
	go func() {
		errc <- srv.ListenAndServe()
	}()

	lg.Info("received signal, shutting down server", zap.String("signal", (<-notifier).String()))
	rootCancel()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	srv.Shutdown(ctx)
	cancel()
	lg.Info("shut down server", zap.Error(<-errc))

	signal.Stop(notifier)
}
