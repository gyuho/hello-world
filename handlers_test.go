package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gyuho/hello-world/version"

	"go.uber.org/zap"
)

func TestHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/hello-world", &contextAdapter{
		lg:      zap.NewExample(),
		ctx:     context.Background(),
		handler: contextHandlerFunc(helloWorldHandler),
	})
	mux.Handle("/readiness", &contextAdapter{
		lg:      zap.NewExample(),
		ctx:     context.Background(),
		handler: contextHandlerFunc(readinessHandler),
	})
	mux.Handle("/liveness", &contextAdapter{
		lg:      zap.NewExample(),
		ctx:     context.Background(),
		handler: contextHandlerFunc(livenessHandler),
	})
	mux.Handle("/status", &contextAdapter{
		lg:      zap.NewExample(),
		ctx:     context.Background(),
		handler: contextHandlerFunc(statusHandler),
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	rs, err := http.Get(ts.URL + "/hello-world")
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(body, []byte("<b>Hello World!</b>")) {
		t.Fatalf("expected %q, got %q", "<b>Hello World!</b>", string(body))
	}

	rs, err = http.Get(ts.URL + "/readiness")
	if err != nil {
		t.Fatal(err)
	}
	if rs.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, rs.StatusCode)
	}

	rs, err = http.Post(ts.URL+"/readiness", "text/plain", nil)
	if err != nil {
		t.Fatal(err)
	}
	if rs.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status code %v, got %v", http.StatusMethodNotAllowed, rs.StatusCode)
	}

	rs, err = http.Get(ts.URL + "/readiness")
	if err != nil {
		t.Fatal(err)
	}
	if rs.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, rs.StatusCode)
	}

	rs, err = http.Get(ts.URL + "/liveness")
	if err != nil {
		t.Fatal(err)
	}
	if rs.StatusCode != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, rs.StatusCode)
	}
	body, err = ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	rs.Body.Close()
	if !bytes.Equal(body, []byte("LIVE\n")) {
		t.Errorf("expected 'LIVE', got %q", string(body))
	}

	rs, err = http.Get(ts.URL + "/not-exists")
	if err != nil {
		t.Fatal(err)
	}
	if rs.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %v, got %v", http.StatusNotFound, rs.StatusCode)
	}

	version.GitCommit, version.ReleaseVersion, version.BuildTime = "a", "b", "c"
	rs, err = http.Get(ts.URL + "/status")
	if err != nil {
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	rs.Body.Close()
	var v version.Version
	if err = json.Unmarshal(body, &v); err != nil {
		t.Fatal(err)
	}
	if v.GitCommit != version.GitCommit {
		t.Errorf("GitCommit expected %q, got %q", version.GitCommit, v.GitCommit)
	}
	if v.ReleaseVersion != version.ReleaseVersion {
		t.Errorf("ReleaseVersion expected %q, got %q", version.ReleaseVersion, v.ReleaseVersion)
	}
	if v.BuildTime != version.BuildTime {
		t.Errorf("BuildTime expected %q, got %q", version.BuildTime, v.BuildTime)
	}
}
