package request

import (
	"net/http"
	"sync"
)

type Requester struct {
	client  *http.Client
	request *http.Request
	mutex   *sync.Mutex
}
