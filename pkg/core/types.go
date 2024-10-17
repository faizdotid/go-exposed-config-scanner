package core

import (
	"go-exposed-config-scanner/pkg/matcher"
	"go-exposed-config-scanner/pkg/request"
	"sync"
	"sync/atomic"
)

type Scanner struct {
	name       string
	matchFrom  string
	matcher    matcher.IMatcher
	client     *request.Requester
	output     string
	verbose    bool
	matchOnly  bool
	counter    *atomic.Uint64
	totalCount uint64
	mu         *sync.Mutex
}
