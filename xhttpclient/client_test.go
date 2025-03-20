package xhttpclient

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type TestStruct struct {
	Test1 string `json:"test1"`
	Test2 string `json:"test2"`
}

type TestLogger struct{}

func (l *TestLogger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func TestTrackProgress(t *testing.T) {
	var data []byte
	client := New("test", WithTimeout(240*time.Second), WithTrackProgress(true), WithLog(t.Logf))
	statusCode, err := client.Do(context.Background(), "GET", "https://ash-speed.hetzner.com/100MB.bin", nil, &data)
	if err != nil {
		t.Fatalf("fatal err: %v", err.Error())
	}

	t.Logf("statusCode: %v", statusCode)
}

func TestClient(t *testing.T) {
	retries := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request: %v", r.URL)
		if retries < 2 {
			retries++
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if r.URL.Path != "/test" {
			t.Errorf("Expected to request '/test', got: %s", r.URL.Path)
		}

		values := r.URL.Query()

		t.Logf("query params: %v", values)

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}

		if r.Body != nil {
			data, _ := io.ReadAll(r.Body)
			t.Logf("body: %v", string(data))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"test1": "test1", "test2": "test2"}`))
	}))

	defer server.Close()

	var ts TestStruct
	var logger = &TestLogger{}

	ts.Test1 = "ts1"
	ts.Test2 = "ts2"

	client := New(server.URL, WithLog(logger.Debugf))
	statusCode, err := client.Do(context.Background(), "GET", "/test", &ts, &ts, url.Values{"test": {"test1"}, "test2": {"test2"}})
	if err != nil {
		log.Printf("err statusCode: %d", statusCode)
		log.Printf("context cancelled? %v", errors.Is(err, context.Canceled))
		log.Printf("context deadline exceeded? %v", errors.Is(err, context.DeadlineExceeded))
		log.Printf("err: %v", err)
	} else {
		log.Printf("statusCode: %d", statusCode)
		log.Printf("test struct: %v", ts)
	}
}
