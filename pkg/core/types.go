package core

import (
	"go-exposed-config-scanner/pkg/matcher"
	"go-exposed-config-scanner/pkg/request"
	"os"
)

type Scanner struct {
	client    *request.Requester
	matcher   matcher.Matcher
	output    *os.File
	name      string
	matchFrom string
	verbose   bool
	matchOnly bool
}
