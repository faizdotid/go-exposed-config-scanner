package request

import (
	"crypto/tls"
	"go-exposed-config-scanner/pkg/templates"
	"net/http"
	"net/url"
)

func NewRequester(req templates.Request) (*Requester, error) {
	httpReq, err := http.NewRequest(req.Method, "", req.Body)
	if err != nil {
		return nil, err
	}

	httpReq.Header = req.Headers

	r := &Requester{
		client: &http.Client{
			Timeout: req.Timeout,
			Transport: &http.Transport{
				MaxIdleConns: 10,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // skip certificate verification
				},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error { // prevent redirect
				return http.ErrUseLastResponse
			},
		},
		request: httpReq,
	}

	return r, nil
}
func (r *Requester) newRequest(target string) (*http.Request, error) {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	return &http.Request{
		Method: r.request.Method,
		URL:    parsedURL,
		Header: r.request.Header,
		Body:   r.request.Body,
	}, nil

}
func (r *Requester) Do(target string) (*http.Response, error) {
	copyReq, err := r.newRequest(target)
	if err != nil {
		return nil, err
	}

	return r.client.Do(copyReq)
}
