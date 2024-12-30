package main

import (
	"context"
	"fmt"
	"go-exposed-config-scanner/internal/cli"
	"go-exposed-config-scanner/internal/helpers"
	"go-exposed-config-scanner/internal/utils"
	"go-exposed-config-scanner/pkg/color"
	"go-exposed-config-scanner/pkg/core"
	"go-exposed-config-scanner/pkg/request"
	"go-exposed-config-scanner/pkg/templates"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	log.SetOutput(io.Discard)
}

func main() {
	ctx := context.Background()
	scanStartTime := time.Now()

	appTemplates := templates.Templates{}
	if err := appTemplates.LoadTemplate("templates"); err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}

	args := cli.ParseArgs(&appTemplates)
	targetTemplates, err := helpers.ParseArgsForTemplates(args.TemplateId, args.All, &appTemplates)
	if err != nil {
		log.Fatalf("failed to filter templates: %v", err)
	}

	if _, err := os.Stat("results"); os.IsNotExist(err) {
		os.Mkdir("results", 0755)
	}

	fileContent, err := os.ReadFile(args.List)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	targetURLs := strings.Split(strings.TrimSpace(string(fileContent)), "\n")

	var totalURLCount uint64
	for _, template := range targetTemplates {
		totalURLCount += uint64(len(targetURLs)) * uint64(len(template.Paths))
	}

	fmt.Printf("%s %s %s\n",
		color.White.AnsiFormat("[")+color.Cyan.AnsiFormat("INFO")+color.White.AnsiFormat("]"),
		color.Green.AnsiFormat(fmt.Sprintf("Loaded %d URLs.", totalURLCount)),
		color.Yellow.AnsiFormat("Starting scan..."))

	var scannedURLCount atomic.Uint64
	var wg sync.WaitGroup
	var semaphore = make(chan struct{}, args.Threads)

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	wg.Add(len(targetTemplates))

	maxIdleConns := args.Threads / len(targetTemplates)
	for _, template := range targetTemplates {
		go func(t *templates.Template) {
			defer wg.Done()
			initializeScanner(
				ctx,
				&utils.ScanOptions{
					Template:        t,
					TargetURLs:      targetURLs,
					TotalURLCount:   totalURLCount,
					ScannedURLCount: &scannedURLCount,
					CliArgs:         args,
					Semaphore:       semaphore,
				},
				maxIdleConns,
			)
		}(template)
	}

	wg.Wait()

	scanDuration := time.Since(scanStartTime)
	fmt.Printf("%s %s %s\n",
		color.White.AnsiFormat("[")+color.Cyan.AnsiFormat("INFO")+color.White.AnsiFormat("]"),
		color.Green.AnsiFormat("Scan completed."),
		color.Yellow.AnsiFormat(fmt.Sprintf("Elapsed time: %s", scanDuration)))
}

func initializeScanner(ctx context.Context, opts *utils.ScanOptions, maxIdleConns int) {
	if opts.CliArgs.Timeout == 0 && opts.Template.Request.Timeout == 0 {
		opts.Template.Request.Timeout = time.Duration(opts.CliArgs.Timeout) * time.Second
	}

	requester, err := request.NewRequester(*opts.Template.Request, maxIdleConns)
	if err != nil {
		log.Fatalf("Failed to create requester: %v", err)
	}

	outputFile := fmt.Sprintf("results/%s", opts.Template.Output)

	scanner := core.NewScanner(
		requester,
		opts.Template.Matcher,
		outputFile,
		opts.Template.Name,
		opts.Template.MatchFrom,
		opts.CliArgs.Verbose,
		opts.CliArgs.MatchOnly,
		opts.ScannedURLCount,
		opts.TotalURLCount,
	)

	targetURLChannel := make(chan string, opts.CliArgs.Threads)

	go func() {
		helpers.MergeURLAndPaths(opts.TargetURLs, opts.Template.Paths, targetURLChannel)
		close(targetURLChannel)
	}()

	var wg sync.WaitGroup

	for targetURL := range targetURLChannel {
		opts.Semaphore <- struct{}{}
		wg.Add(1)
		go func(url string) {
			scanner.Scan(ctx, url, &wg)
			<-opts.Semaphore
		}(targetURL)
	}

	wg.Wait()
}
