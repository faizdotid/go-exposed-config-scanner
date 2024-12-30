package utils

import (
	"go-exposed-config-scanner/internal/cli"
	"go-exposed-config-scanner/pkg/templates"
	"sync/atomic"
)

type ScanOptions struct {
	Template        *templates.Template
	TargetURLs      []string
	TotalURLCount   uint64
	ScannedURLCount *atomic.Uint64
	CliArgs         *cli.Args
	Semaphore       chan struct{}
}
