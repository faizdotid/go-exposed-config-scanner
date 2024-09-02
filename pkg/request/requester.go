package request

import (
	"crypto/tls"
	"go-exposed-config-scanner/pkg/templates"
	"net/http"
	"sync"
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
		},
		request: httpReq,
		mutex:   &sync.Mutex{},
	}

	return r, nil
}

func (r *Requester) Do(target string) (*http.Response, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.setURLRequest(target); err != nil {
		return nil, err
	}

	return r.client.Do(r.request)
}
