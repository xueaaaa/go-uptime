package worker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/xueaaaa/go-uptime/internal/config"
	"github.com/xueaaaa/go-uptime/internal/job"
	"github.com/xueaaaa/go-uptime/internal/site/model"
)

var internalErrorMsg = "internal error"

func statusOKHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func statusOKExceedTimeoutHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	w.WriteHeader(http.StatusOK)
}

func statusInternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(internalErrorMsg))
}

func TestSiteResponds(t *testing.T) {
	tests := []struct {
		name               string
		handler            http.HandlerFunc
		expectedStatus     model.Status
		expectedStatusCode int
		expectedError      string
	}{
		{
			name:               "status OK",
			handler:            statusOKHandler,
			expectedStatus:     model.Available,
			expectedStatusCode: http.StatusOK,
			expectedError:      "",
		},
		{
			name:               "status Internal server error",
			handler:            statusInternalServerErrorHandler,
			expectedStatus:     model.Unavailable,
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      internalErrorMsg,
		},
		{
			name:               "status OK exceed timeout",
			handler:            statusOKExceedTimeoutHandler,
			expectedStatus:     model.Unavailable,
			expectedStatusCode: 0,
			expectedError:      "context deadline exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			pool := &Pool{}
			result := pool.check(
				context.Background(),
				job.Job{URL: server.URL},
				config.Config{TimeoutAfter: time.Second},
			)

			if result.Status != tt.expectedStatus {
				t.Fatalf("expected %v, got %v", tt.expectedStatus, result.Status)
			}
			if result.StatusCode != tt.expectedStatusCode {
				t.Fatalf("expected %d, got %d", tt.expectedStatusCode, result.StatusCode)
			}
			if !strings.Contains(result.Error, tt.expectedError) {
				t.Fatalf("expected error %q, got %q", tt.expectedError, result.Error)
			}
		})
	}
}
