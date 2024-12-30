package request_test

import (
	"context"
	"go-exposed-config-scanner/pkg/request"
	"go-exposed-config-scanner/pkg/templates"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	testTimeout    = 5 * time.Second
	maxIdleConns   = 100
	testBaseURL    = "http://httpbin.org"
	sampleFilesURL = "https://getsamplefiles.com/download/zip"
)

func setupTestContext(t *testing.T) (context.Context, context.CancelFunc) {
	_ = t
	return context.WithTimeout(context.Background(), testTimeout)
}

func TestHTTPMethods(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		url           string
		headers       http.Header
		body          string
		expectedCode  int
		checkResponse func(*testing.T, *http.Response)
	}{
		{
			name:         "GET Request",
			method:       "GET",
			url:          testBaseURL + "/get",
			expectedCode: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Failed to read response body: %v", err)
				}
				t.Logf("Response body: %s", body)
			},
		},
		{
			name:   "POST Request",
			method: "POST",
			url:    testBaseURL + "/post",
			body:   `{"key": "value"}`,
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Headers Test",
			method:       "HEAD",
			url:          testBaseURL + "/headers",
			expectedCode: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				for k, v := range resp.Header {
					t.Logf("Header %s: %v", k, v)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := setupTestContext(t)
			defer cancel()

			var bodyReader io.ReadCloser
			if tt.body != "" {
				bodyReader = io.NopCloser(strings.NewReader(tt.body))
			}

			req := templates.Request{
				Method:  tt.method,
				Timeout: testTimeout,
				Headers: tt.headers,
				Body:    bodyReader,
			}

			r, err := request.NewRequester(req, maxIdleConns)
			if err != nil {
				t.Fatalf("Failed to create requester: %v", err)
			}

			resp, err := r.Do(ctx, tt.url)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}

func TestTemplateRequests(t *testing.T) {
	ctx, cancel := setupTestContext(t)
	defer cancel()

	var templates templates.Templates
	if err := templates.LoadTemplate("../../templates"); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	temp := templates[0]
	client, err := request.NewRequester(*temp.Request, maxIdleConns)
	if err != nil {
		t.Fatalf("Failed to create requester: %v", err)
	}

	sampleFiles := []string{
		"sample-1.zip",
		"sample-2.zip",
		"sample-3.zip",
		"sample-4.zip",
		"sample-5.zip",
	}

	for _, file := range sampleFiles {
		t.Run("Testing "+file, func(t *testing.T) {
			resp, err := client.Do(ctx, sampleFilesURL+"/"+file)
			if err != nil {
				t.Fatalf("Request failed for %s: %v", file, err)
			}
			defer resp.Body.Close()

			match, err := temp.Matcher.Match(resp)
			if err != nil {
				t.Errorf("Matcher failed for %s: %v", file, err)
			}

			t.Logf("File: %s, Matched: %v, Headers: %v", file, match, resp.Header)
		})
	}
}

func BenchmarkRequests(b *testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var templates templates.Templates
	if err := templates.LoadTemplate("../../templates"); err != nil {
		b.Fatalf("Failed to load templates: %v", err)
	}

	client, err := request.NewRequester(*templates[0].Request, maxIdleConns)
	if err != nil {
		b.Fatalf("Failed to create requester: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.Do(ctx, sampleFilesURL+"/sample-1.zip")
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
		resp.Body.Close()
	}
}
