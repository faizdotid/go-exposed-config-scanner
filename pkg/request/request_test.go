package request_test

import (
	"fmt"
	"go-exposed-config-scanner/pkg/request"
	"go-exposed-config-scanner/pkg/templates"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRequest(t *testing.T) {
	rawJson := `{"key": "value"}`
	body := io.NopCloser(strings.NewReader(rawJson))
	t.Run("TestGET", func(t *testing.T) {
		req := templates.Request{
			Method:  "GET",
			Timeout: 5,
			Headers: http.Header{},
			Body:    nil,
		}
		r, err := request.NewRequester(req)
		if err != nil {
			t.Errorf("Error creating requester: %v", err)
		}
		resp, err := r.Do("http://example.com")
		if err != nil {
			t.Errorf("Error making request: %v", err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Error reading response body: %v", err)
		}

		t.Logf("Response body: %s", body)
	})
	t.Run("TestPOST", func(t *testing.T) {
		req := templates.Request{
			Method:  "POST",
			Timeout: 5,
			Headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: body,
		}
		r, err := request.NewRequester(req)
		if err != nil {
			t.Errorf("Error creating requester: %v", err)
		}
		resp, err := r.Do("http://httpbin.org/post")
		if err != nil {
			t.Errorf("Error making request: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("Expected status code 200, got %d", resp.StatusCode)
		}
	})
	t.Run("TestHeaders", func(t *testing.T) {
		req := templates.Request{
			Method:  "HEAD",
			Timeout: 5,
			Headers: http.Header{},
			Body:    nil,
		}
		r, err := request.NewRequester(req)
		if err != nil {
			t.Errorf("Error creating requester: %v", err)
		}
		resp, err := r.Do("https://github.com/faizdotid/go-simple-web-server/archive/refs/heads/main.zip")

		if err != nil {
			t.Errorf("Error making request: %v", err)
		}
		for k, v := range resp.Header {
			t.Logf("%s: %v", k, v)
		}
		fmt.Printf("%v", req.Headers)
	})
}

func TestRequestWithTemplate(t *testing.T) {
	var mytemp templates.Templates
	mytemp.LoadTemplate("../../configs")
	temp := mytemp[0]
	client, err := request.NewRequester(*temp.Request)
	if err != nil {
		t.Errorf("Error creating requester: %v", err)
	}
	temp.Paths = append(temp.Paths, "sample-5.zip")
	target := "https://getsamplefiles.com/download/zip/"
	for _, path := range []string{"sample-1.zip", "sample-2.zip", "sample-3.zip", "sample-4.zip", "sample-5.zip"} {
		resp, err := client.Do(target + path)
		if err != nil {
			t.Errorf("Error making request: %v", err)
		}

		headers := fmt.Sprintf("%v", resp.Header)
		if temp.Match.Match([]byte(headers)) {
			t.Logf("URL: %s", resp.Request.URL.String())
			t.Logf("Matched: %s", headers)
		}

	}
}

func BenchmarkTestRequestWithTemplate(b *testing.B) {
	var mytemp templates.Templates
	mytemp.LoadTemplate("../../configs")
	temp := mytemp[0]
	// temp.Request.Method = "GET"
	client, err := request.NewRequester(*temp.Request)
	_ = client
	if err != nil {
		b.Errorf("Error creating requester: %v", err)
	}
	for i := 0; i < b.N; i++ {
		client.Do("https://getsamplefiles.com/download/zip/sample-5.zip")
	}
}
