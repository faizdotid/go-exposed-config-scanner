package core_test

import (
	"crypto/tls"
	"go-exposed-config-scanner/pkg/core"
	// "go-exposed-config-scanner/pkg/core"
	"net/http"
	"testing"
)

func BenchmarkCore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		httpClient := &http.Client{
			Timeout: 5,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		core.NewScanner(httpClient, nil, nil, "", "")
		
	}
}