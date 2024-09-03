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

	if _, err := os.Stat("results"); os.IsNotExist(err) {
		os.Mkdir("results", 0755)
	}

	// read file
	fileContent, err := os.ReadFile(args.FileList)
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
	time.Sleep(1 * time.Second)

	// initialize a counter for URLs
	urlCount := atomic.Uint64{}
	urlCount.Store(1)

	// initialize mutex and wait group
	var mu sync.Mutex
	var wg sync.WaitGroup

	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize the scanner for each selected template
	wg.Add(len(selectedTemplates))
	for _, template := range selectedTemplates {
		go func(t *templates.Template) {
			defer wg.Done()
			initializeScanner(t, urls, totalCount, &urlCount, &mu, args)
		}(template)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Printf("%s %s %s\n",
		color.White.AnsiFormat("[")+color.Cyan.AnsiFormat("INFO")+color.White.AnsiFormat("]"),
		color.Green.AnsiFormat("Scan completed."),
		color.Yellow.AnsiFormat(fmt.Sprintf("Elapsed time: %s", elapsedTime)))
}

// initializeScanner sets up and runs the scan for a specific template using multiple threads
func initializeScanner(template *templates.Template, urls []string, totalCount uint64, urlCount *atomic.Uint64, mu *sync.Mutex, args *cli.Args) {
	if args.Timeout > 0 {
		template.Request.Timeout = time.Duration(args.Timeout) * time.Second
	}

	requester, err := request.NewRequester(*template.Request)
	if err != nil {
		log.Fatalf("failed to create requester: %v", err)
	}

	fileOutput, err := os.OpenFile(fmt.Sprintf("results/%s", template.Output), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer fileOutput.Close()

	scanner := core.NewScanner(requester, template.Match, fileOutput, template.Name, template.MatchFrom, args.Verbose, args.MatchOnly)

	threadsChannel := make(chan struct{}, args.Threads)
	var wg sync.WaitGroup

	var targets []string
	for _, url := range urls {
		helpers.MergeURLAndPaths(url, template.Paths, &targets)
	}

	for _, target := range targets {
		threadsChannel <- struct{}{}
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			defer func() { <-threadsChannel }()

			mu.Lock()
			urlCount.Add(1)
			mu.Unlock()
			scanner.Scan(t, urlCount.Load(), totalCount)
		}(target)
	}

	wg.Wait()
}
