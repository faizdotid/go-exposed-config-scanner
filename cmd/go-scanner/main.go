package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"go-exposed-config-scanner/internal/helpers"
	"go-exposed-config-scanner/pkg/color"
	"go-exposed-config-scanner/pkg/core"
	"go-exposed-config-scanner/pkg/templates"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// Parse command-line flags
	templateID := flag.String("id", "", "Template ID, comma-separated for multiple templates")
	scanAll := flag.Bool("all", false, "Scan all templates")
	listFile := flag.String("list", "", "List of URLs to scan")
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
		const columnWidth = 30

		// Define colored headers
		// Function to center text within a column
		centerText := func(text string, width int, colors ...color.Color) string {
			padding := (width - len(text)) / 2
			currrentStr := strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
			return color.Coloring(currrentStr, colors...)
		}

		// Print the centered header with separators
		fmt.Printf("|%-5s|%-*s|%-*s|%-*s|%-*s|\n",
			centerText("No.", 5, color.Red, color.Bold),
			columnWidth, centerText("ID", columnWidth, color.Red, color.Bold),
			columnWidth, centerText("Name", columnWidth, color.Red, color.Bold),
			columnWidth, centerText("Match From", columnWidth, color.Red, color.Bold),
			columnWidth, centerText("Paths", columnWidth, color.Red, color.Bold),
		)

		fmt.Println(strings.Repeat("-", (columnWidth*4+4)+6) + "|")
		for idx, iter := range templateList {
			// Join paths into a single string

			var joinedPaths string = strings.Join(iter.Paths, ", ")
			// Truncate paths if they exceed the column width
			if len(joinedPaths) > columnWidth {
				joinedPaths = joinedPaths[:columnWidth-5] + "..."
			}

			// Format template fields with colors

			// Print the centered template details with separators
			fmt.Printf("|%-5s|%-*s|%-*s|%-*s|%-*s|\n",
				centerText(strconv.Itoa(idx+1), 5, color.Blue, color.Bold),
				columnWidth, centerText(iter.ID, columnWidth, color.Green),
				columnWidth, centerText(iter.Name, columnWidth, color.Blue),
				columnWidth, centerText(iter.MatchFrom, columnWidth, color.Blue),
				columnWidth, centerText(joinedPaths, columnWidth, color.Blue),
			)
		}

		os.Exit(0)
	}

	if *listFile == "" {
		fmt.Println("Usage: go-scanner -list <list> [-id <template_id> | -all] [-threads <count>] [-timeout <seconds>]")
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
	fileContent, err := os.ReadFile(*listFile)
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
	lock := sync.Mutex{}
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Start scanning for each template
	for _, template := range selectedTemplates {
		go func(t *templates.Template) {
			defer wg.Done()
			runScanner(t, client, *threadCount, &urlCount, totalCount, urls, &lock)
		}(template)
	}

	// Wait for all scans to complete
	wg.Wait()

	log.Printf("Elapsed time: %v", time.Since(startTime))
}

// runScanner executes the scan for a specific template using multiple threads
func runScanner(t *templates.Template, client *http.Client, threads int, urlCount *atomic.Uint64, totalCount uint64, urls []string, lock *sync.Mutex) {
	// Open or create the output file
	outputFile, err := os.OpenFile(fmt.Sprintf("results/%s", t.Output), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	defer outputFile.Close()

	// Initialize the scanner with the HTTP client, validator, and output file
	myscanner := core.NewScanner(client, t.Matcher, outputFile, t.Name, t.MatchFrom)

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
			defer func() { <-threadLimiter }() // Release the thread limiter

			lock.Lock()
			urlCount.Add(1)
			lock.Unlock()
			myscanner.Scan(url, urlCount.Load(), totalCount)
		}(target)
	}

	wg.Wait()

}
