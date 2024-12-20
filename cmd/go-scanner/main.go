package main

import (
	"fmt"
	"go-exposed-config-scanner/internal/cli"
	"go-exposed-config-scanner/internal/helpers"
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

type ScanOptions struct {
	template        *templates.Template
	targetURLs      []string
	totalURLCount   uint64
	scannedURLCount *atomic.Uint64
	mutex           *sync.Mutex
	cliArgs         *cli.Args
	threadLimiter   chan struct{}
}

func main() {
	scanStartTime := time.Now()

	// set log into discard
	log.SetOutput(io.Discard)

	// load available templates
	availableTemplates := templates.Templates{}
	if err := availableTemplates.LoadTemplate("templates"); err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	// parse command line arguments
	cliArgs := cli.ParseArgs(availableTemplates)

	// filtering templates based on the provided arguments
	selectedTemplates, err := helpers.ParseArgsForTemplates(cliArgs.TemplateId, cliArgs.All, &availableTemplates)
	if err != nil {
		log.Fatalf("Failed to parse args for templates: %v", err)
	}

	// create results directory if not exists
	if _, err := os.Stat("results"); os.IsNotExist(err) {
		os.Mkdir("results", 0755)
	}

	// read target urls from file
	fileContent, err := os.ReadFile(cliArgs.List)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	targetURLs := strings.Split(strings.TrimSpace(string(fileContent)), "\n")

	// calculate total url to be scanned
	var totalURLCount uint64
	for _, template := range selectedTemplates {
		totalURLCount += uint64(len(targetURLs)) * uint64(len(template.Paths))
	}

	fmt.Printf("%s %s %s\n",
		color.White.AnsiFormat("[")+color.Cyan.AnsiFormat("INFO")+color.White.AnsiFormat("]"),
		color.Green.AnsiFormat(fmt.Sprintf("Loaded %d URLs.", totalURLCount)),
		color.Yellow.AnsiFormat("Starting scan..."))

	var scannedURLCount atomic.Uint64
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var threadLimiter = make(chan struct{}, cliArgs.Threads)

	runtime.GOMAXPROCS(runtime.NumCPU())

	// initialize scanner for each template
	wg.Add(len(selectedTemplates))
	for _, template := range selectedTemplates {
		go func(t *templates.Template) {
			defer wg.Done()
			initializeScanner(
				&ScanOptions{
					template:        t,
					targetURLs:      targetURLs,
					totalURLCount:   totalURLCount,
					scannedURLCount: &scannedURLCount,
					mutex:           &mutex,
					cliArgs:         cliArgs,
					threadLimiter:   threadLimiter,
				},
			)
		}(template)
	}

	// wait all goroutines to finish
	wg.Wait()

	scanDuration := time.Since(scanStartTime)
	fmt.Printf("%s %s %s\n",
		color.White.AnsiFormat("[")+color.Cyan.AnsiFormat("INFO")+color.White.AnsiFormat("]"),
		color.Green.AnsiFormat("Scan completed."),
		color.Yellow.AnsiFormat(fmt.Sprintf("Elapsed time: %s", scanDuration)))
}

// initializeScanner sets up and runs the scan for a specific template using multiple threads
func initializeScanner(opts *ScanOptions) {
	if opts.cliArgs.Timeout == 0 && opts.template.Request.Timeout == 0 {
		opts.template.Request.Timeout = time.Duration(opts.cliArgs.Timeout) * time.Second
	}

	requester, err := request.NewRequester(*opts.template.Request)
	if err != nil {
		log.Fatalf("Failed to create requester: %v", err)
	}

	outputFile := fmt.Sprintf("results/%s", opts.template.Output)

	scanner := core.NewScanner(
		requester,
		opts.template.Matcher,
		outputFile,
		opts.template.Name,
		opts.template.MatchFrom,
		opts.cliArgs.Verbose,
		opts.cliArgs.MatchOnly,
		opts.scannedURLCount,
		opts.totalURLCount,
		opts.mutex,
	)

	targetURLChannel := make(chan string, opts.cliArgs.Threads)

	go func() {
		helpers.MergeURLAndPaths(opts.targetURLs, opts.template.Paths, targetURLChannel)
		close(targetURLChannel)
	}()

	var wg sync.WaitGroup

	for targetURL := range targetURLChannel {
		opts.threadLimiter <- struct{}{} // acquire the thread limiter
		wg.Add(1)
		go func(url string) {
			scanner.Scan(url, &wg)
			<-opts.threadLimiter // release the thread limiter
		}(targetURL)
	}

	wg.Wait()
}
