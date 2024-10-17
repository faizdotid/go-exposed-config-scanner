package main

import (
	"fmt"
	"go-exposed-config-scanner/internal/cli"
	"go-exposed-config-scanner/internal/helpers"
	"go-exposed-config-scanner/pkg/color"
	"go-exposed-config-scanner/pkg/core"
	"go-exposed-config-scanner/pkg/request"
	"go-exposed-config-scanner/pkg/templates"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Options struct {
	template   *templates.Template
	urls       []string
	totalCount uint64
	urlCount   *atomic.Uint64
	mu         *sync.Mutex
	args       *cli.Args
}

func main() {
	startTime := time.Now()

	// load templates
	currentTemplates := templates.Templates{}
	if err := currentTemplates.LoadTemplate("templates"); err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}

	// parsing command line arguments
	args := cli.ParseArgs(currentTemplates)

	// filter templates
	selectedTemplates, err := helpers.ParseArgsForTemplates(args.TemplateId, args.All, &currentTemplates)
	if err != nil {
		log.Fatalf("failed to parse args for templates: %v", err)
	}

	// create results directory if it doesn't exist
	if _, err := os.Stat("results"); os.IsNotExist(err) {
		os.Mkdir("results", 0755)
	}

	// read file
	fileContent, err := os.ReadFile(args.List)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	urls := strings.Split(strings.TrimSpace(string(fileContent)), "\n")

	// calculate total count of URLs to scan
	var totalCount uint64
	for _, t := range selectedTemplates {
		totalCount += uint64(len(urls)) * uint64(len(t.Paths))
	}

	fmt.Printf("%s %s %s\n",
		color.White.AnsiFormat("[")+color.Cyan.AnsiFormat("INFO")+color.White.AnsiFormat("]"),
		color.Green.AnsiFormat(fmt.Sprintf("Loaded %d URLs.", totalCount)),
		color.Yellow.AnsiFormat("Starting scan..."))

	var urlCount atomic.Uint64
	var mu sync.Mutex
	var wg sync.WaitGroup

	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize the scanner for each selected template
	wg.Add(len(selectedTemplates))
	for _, template := range selectedTemplates {
		go func(t *templates.Template) {
			defer wg.Done()
			initializeScanner(
				&Options{
					template:   t,
					urls:       urls,
					totalCount: totalCount,
					urlCount:   &urlCount,
					mu:         &mu,
					args:       args,
				},
			)
		}(template)
	}

	// wait all goroutines to finish
	wg.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Printf("%s %s %s\n",
		color.White.AnsiFormat("[")+color.Cyan.AnsiFormat("INFO")+color.White.AnsiFormat("]"),
		color.Green.AnsiFormat("Scan completed."),
		color.Yellow.AnsiFormat(fmt.Sprintf("Elapsed time: %s", elapsedTime)))
}

// initializeScanner sets up and runs the scan for a specific template using multiple threads
func initializeScanner(opts *Options) {
	if opts.args.Timeout == 0 && opts.template.Request.Timeout == 0 {
		opts.template.Request.Timeout = time.Duration(opts.args.Timeout) * time.Second
	}

	requester, err := request.NewRequester(*opts.template.Request)
	if err != nil {
		log.Fatalf("failed to create requester: %v", err)
	}

	fileOutput, err := os.OpenFile(fmt.Sprintf("results/%s", opts.template.Output), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer fileOutput.Close()

	scanner := core.NewScanner(
		requester,
		opts.template.Matcher,
		fileOutput,
		opts.template.Name,
		opts.template.MatchFrom,
		opts.args.Verbose,
		opts.args.MatchOnly,
		opts.urlCount,
		opts.totalCount,
		opts.mu,
	)

	threadsChannel := make(chan struct{}, opts.args.Threads)
	targetsChannel := make(chan string, len(opts.urls)*len(opts.template.Paths))

	go func() {
		helpers.MergeURLAndPaths(opts.urls, opts.template.Paths, targetsChannel)
		close(targetsChannel)
	}()

	var wg sync.WaitGroup

	for target := range targetsChannel {
		threadsChannel <- struct{}{}
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			defer func() { <-threadsChannel }() // release thread
			scanner.Scan(t)
		}(target)
	}

	wg.Wait()
}
