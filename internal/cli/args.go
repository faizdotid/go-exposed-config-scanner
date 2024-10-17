package cli

import (
	"flag"
	"fmt"
	"go-exposed-config-scanner/pkg/color"
	"go-exposed-config-scanner/pkg/templates"
	"os"
	"strconv"
	"strings"
)

func init() {
	flag.StringVar(&currentArgs.TemplateId, "id", "", "Template ID, comma-separated for multiple templates")
	flag.BoolVar(&currentArgs.All, "all", false, "Scan all templates")
	flag.StringVar(&currentArgs.List, "list", "", "Path list of urls")
	flag.IntVar(&currentArgs.Threads, "threads", 10, "Number of threads to use")
	flag.BoolVar(&currentArgs.Show, "show", false, "Show available templates")
	flag.BoolVar(&currentArgs.MatchOnly, "match", false, "Print only match URLs")
	flag.BoolVar(&currentArgs.Verbose, "verbose", false, "Print errors verbose")
	flag.IntVar(&currentArgs.Timeout, "timeout", 0, "Timeout for HTTP requests (It will be applied to all templates)")

}

func ParseArgs(templates templates.Templates) *Args {
	flag.Parse()

	if currentArgs.Show {
		ShowTemplates(templates)
		os.Exit(0)
	}

	if currentArgs.List == "" {
		printUsageAndExit()
	}

	return &currentArgs
}

func printUsageAndExit() {
	fmt.Println("Usage: go-scanner -list <list> [-templateId <template_id> | -all] [-threads <count>] [-show]")
	flag.PrintDefaults()
	os.Exit(1)
}

func ShowTemplates(templateList templates.Templates) {
	if len(templateList) == 0 {
		fmt.Println(color.Coloring("No templates found", color.Red, color.Bold))
		return
	}
	const columnWidth = 30

	// Center text within a column and apply colors
	centerText := func(text string, width int, colors ...color.Color) string {
		padding := (width - len(text)) / 2
		currentStr := strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
		return color.Coloring(currentStr, colors...)
	}

	// Print the centered header with separators
	fmt.Printf("|%-5s|%-*s|%-*s|%-*s|%-*s|\n",
		centerText("No.", 5, color.Red, color.Bold),
		columnWidth, centerText("ID", columnWidth, color.Red, color.Bold),
		columnWidth, centerText("Name", columnWidth, color.Red, color.Bold),
		columnWidth, centerText("Output", columnWidth, color.Red, color.Bold),
		columnWidth, centerText("Paths", columnWidth, color.Red, color.Bold),
	)

	fmt.Println(strings.Repeat("-", (columnWidth*4+4)+6) + "|")
	for idx, tmpl := range templateList {
		joinedPaths := truncateString(strings.Join(tmpl.Paths, ", "), columnWidth-5)

		fmt.Printf("|%-5s|%-*s|%-*s|%-*s|%-*s|\n",
			centerText(strconv.Itoa(idx+1), 5, color.Blue, color.Bold),
			columnWidth, centerText(tmpl.ID, columnWidth, color.Green),
			columnWidth, centerText(tmpl.Name, columnWidth, color.Blue),
			columnWidth, centerText(tmpl.Output, columnWidth, color.Blue),
			columnWidth, centerText(joinedPaths, columnWidth, color.Blue),
		)
	}
}

func truncateString(str string, maxLength int) string {
	if len(str) > maxLength {
		return str[:maxLength] + "..."
	}
	return str
}
