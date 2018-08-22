package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gyuho/hello-world/version"

	"go.uber.org/zap"
)

// contextHandler handles ServeHTTP with context.
type contextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request) error
}

// contextHandlerFunc defines HandlerFunc function signature to wrap context.
type contextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request) error

// ServeHTTPContext serve HTTP requests with context.
func (f contextHandlerFunc) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	return f(ctx, w, req)
}

// contextAdapter wraps context handler.
type contextAdapter struct {
	lg      *zap.Logger
	ctx     context.Context
	handler contextHandler
}

func (ca *contextAdapter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := ca.handler.ServeHTTPContext(ca.ctx, w, req); err != nil {
		ca.lg.Warn("failed to serve", zap.String("method", req.Method), zap.String("path", req.URL.Path), zap.Error(err))
	}
}

type key int

const userKey key = 0

func with(h contextHandler) contextHandler {
	return contextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
		userID := ""
		ctx = context.WithValue(ctx, userKey, &userID)
		return h.ServeHTTPContext(ctx, w, req)
	})
}

func helloWorldHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case http.MethodGet:
		w.Write([]byte("<b>Hello World!</b>"))
		return nil

	default:
		http.Error(w, "Method Not Allowed", 405)
		return fmt.Errorf("Method %q Not Allowed", req.Method)
	}
}

func readinessHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		return nil

	default:
		http.Error(w, "Method Not Allowed", 405)
		return fmt.Errorf("Method %q Not Allowed", req.Method)
	}
}

func livenessHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("LIVE\n"))
		return err

	default:
		http.Error(w, "Method Not Allowed", 405)
		return fmt.Errorf("Method %q Not Allowed", req.Method)
	}
}

func statusHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case http.MethodGet:
		h, _ := os.Hostname()
		return json.NewEncoder(w).Encode(version.Version{
			GitCommit:      version.GitCommit,
			ReleaseVersion: version.ReleaseVersion,
			BuildTime:      version.BuildTime,
			HostName:       h,
		})

	default:
		http.Error(w, "Method Not Allowed", 405)
		return nil
	}
}
