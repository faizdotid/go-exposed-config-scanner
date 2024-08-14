package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"go-exposed-config-scanner/internal/helpers"
	"go-exposed-config-scanner/pkg/core"
	"go-exposed-config-scanner/pkg/templates"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)


func main() {
	// Parse command-line flags
	templateID := flag.String("id", "", "Template ID, comma-separated for multiple templates")
	scanAll := flag.Bool("all", false, "Scan all templates")
	filePath := flag.String("file", "", "File to scan")
	threadCount := flag.Int("threads", 1, "Number of threads")
	timeout := flag.Int("timeout", 10, "Timeout in seconds")
	show := flag.Bool("show", false, "Show templates")
	flag.Parse()

	// Load templates from the configuration directory
	templateList := templates.Templates{}
	if err := templateList.LoadTemplates("configs"); err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}
	if *show {
		for _, t := range templateList {
			fmt.Printf("ID: %s\nName: %s\nPaths: %v\nOutput: %s\n\n", t.ID, t.Name, t.Paths[0] + " ...", t.Output)
		}
		os.Exit(0)
	}

	if *filePath == "" {
		fmt.Println("Usage: go-scanner -file <file> [-id <template_id> | -all] [-threads <count>] [-timeout <seconds>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create an HTTP client with a custom transport
	client := &http.Client{
		Timeout: time.Duration(*timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}


	// Select the templates based on command-line flags
	selectedTemplates, err := helpers.ParseArgsForTemplates(*templateID, *scanAll, &templateList)
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	// Read the file containing URLs to scan
	fileContent, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	urls := strings.Split(strings.TrimSpace(string(fileContent)), "\n")

	// Calculate the total number of URLs to scan
	totalCount := uint64(0)
	for _, t := range selectedTemplates {
		totalCount += uint64(len(urls) * len(t.Paths))
	}

	// Initialize an atomic counter for the scanned URLs
	var urlCount atomic.Uint64

	urlCount.Store(1)
	// Create a WaitGroup to manage goroutines
	var wg sync.WaitGroup
	wg.Add(len(selectedTemplates))

	fmt.Printf("Starting scan %d URLs with %d threads\n", totalCount, *threadCount)
	startTime := time.Now()

	if _, err := os.Stat("results"); os.IsNotExist(err) {
		os.Mkdir("results", 0755)
	}

	// Start scanning for each template
	for _, template := range selectedTemplates {
		go func(t *templates.Template) {
			defer wg.Done()
			runScanner(t, client, *threadCount, &urlCount, totalCount, urls)
		}(template)
	}

	// Wait for all scans to complete
	wg.Wait()


	log.Printf("Elapsed time: %v", time.Since(startTime))
}

// runScanner executes the scan for a specific template using multiple threads
func runScanner(t *templates.Template, client *http.Client, threads int, urlCount *atomic.Uint64, totalCount uint64, urls []string) {
	// Open or create the output file
	outputFile, err := os.OpenFile(fmt.Sprintf("results/%s", t.Output), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	defer outputFile.Close()

	// Initialize the scanner with the HTTP client, validator, and output file
	myscanner := core.NewScanner(client, t.Validator, outputFile, t.Name, t.MatchFrom)

	// Create a channel to limit the number of concurrent goroutines
	threadLimiter := make(chan struct{}, threads)
	var wg sync.WaitGroup

	var targets []string
	for _, url := range urls {
		helpers.MergeURLAndPaths(url, t.Paths, &targets)
	}

	for _, target := range targets {
		threadLimiter <- struct{}{}
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			defer func() { <-threadLimiter }()
			myscanner.Scan(url, urlCount.Load(), totalCount)
			urlCount.Add(1)
		}(target)
	}
	
	wg.Wait()

}
