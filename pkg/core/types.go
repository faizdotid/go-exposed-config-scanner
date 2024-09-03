package core

import (
	"go-exposed-config-scanner/pkg/matcher"
	"go-exposed-config-scanner/pkg/request"
	"os"
)

type Scanner struct {
	name      string
	matchFrom string
	matcher   matcher.Matcher
	client    *request.Requester
	output    *os.File
	verbose   bool
	matchOnly bool
}
