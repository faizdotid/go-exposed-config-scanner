package request

import (
	"context"
	"crypto/tls"
	"go-exposed-config-scanner/pkg/templates"
	"net/http"
	"net/url"
)

func NewRequester(req templates.Request, maxIdleConns int) (*Requester, error) {
	if maxIdleConns <= 0 {
		maxIdleConns = 100
	}

	httpReq, err := http.NewRequest(req.Method, "", req.Body)
	if err != nil {
		return nil, err
	}

	headers := make(http.Header)
	for k, v := range req.Headers {
		headers[k] = v
	}
	httpReq.Header = headers

	transport := &http.Transport{
		MaxIdleConns: maxIdleConns,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	r := &Requester{
		client: &http.Client{
			Timeout:   req.Timeout,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 3 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
		request: httpReq,
	}

	return r, nil
}

func (r *Requester) newRequest(ctx context.Context, target string) (*http.Request, error) {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	// REWRITE: Clone request using request.Clone
	newReq := r.request.Clone(ctx)
	newReq.URL = parsedURL

	return newReq, nil
}

func (r *Requester) Do(ctx context.Context, target string) (*http.Response, error) {

	req, err := r.newRequest(ctx, target)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req)
}
