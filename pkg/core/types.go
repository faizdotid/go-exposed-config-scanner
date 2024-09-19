package core

import (
	"go-exposed-config-scanner/pkg/matcher"
	"go-exposed-config-scanner/pkg/request"
	"os"
	"sync"
	"sync/atomic"
)

type Scanner struct {
	name       string
	matchFrom  string
	matcher    matcher.Matcher
	client     *request.Requester
	output     *os.File
	verbose    bool
	matchOnly  bool
	counter    *atomic.Uint64
	totalCount uint64
	mu         *sync.Mutex
}
