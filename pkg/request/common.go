package request

import "net/url"

func (r *Requester) setURLRequest(target string) error {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return err
	}
	r.request.URL = parsedURL
	return nil
}
