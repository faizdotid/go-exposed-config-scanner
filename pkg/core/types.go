package core

import (
	"sync/atomic"
	"time"

	"go-exposed-config-scanner/pkg/matcher"
	"go-exposed-config-scanner/pkg/request"
)

const DefaultTimeout = 10 * time.Second

type Scanner struct {
	client     *request.Requester
	matcher    matcher.IMatcher
	name       string
	matchFrom  string
	output     string
	verbose    bool
	matchOnly  bool
	counter    *atomic.Uint64
	totalCount uint64
}
