package xclientv2

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestStruct struct {
	Test1 string `json:"test1"`
	Test2 string `json:"test2"`
}

type TestLogger struct{}

func (l *TestLogger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func TestClient(t *testing.T) {
	retries := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if retries < 2 {
			retries++
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if r.Header.Get("Test") != "test" {
			t.Errorf("Expected Test: test header, got: %s", r.Header.Get("Test"))
		}

		if r.URL.Path != "/test" {
			t.Errorf("Expected to request '/test', got: %s", r.URL.Path)
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"test1": "test1", "test2": "test2"}`))
	}))

	defer server.Close()

	var ts TestStruct
	var logger = &TestLogger{}

	client := New(server.URL, WithLog(logger.Debugf))
	statusCode, err := client.Do(context.Background(), "GET", "test", nil, &ts, http.Header{"test": []string{"test"}})
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
