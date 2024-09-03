package request

import (
	"net/http"
)

type Requester struct {
	client  *http.Client
	request *http.Request
}
