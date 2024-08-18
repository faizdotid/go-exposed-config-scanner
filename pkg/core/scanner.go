package core

import (
	"fmt"
	"go-exposed-config-scanner/pkg/color"
	"go-exposed-config-scanner/pkg/matcher"
	"go-exposed-config-scanner/pkg/utils"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Scanner struct {
	c         *http.Client
	m         matcher.Matcher
	o         *os.File
	name      string
	matchFrom string
}

func NewScanner(client *http.Client, m matcher.Matcher, output *os.File, name string, matchFrom string) *Scanner {
	return &Scanner{
		c:         client,
		m:         m,
		o:         output,
		name:      name,
		matchFrom: matchFrom,
	}
}

func (s *Scanner) Scan(url string, count uint64, totalCount uint64) {
	resp, err := s.c.Get(url)
	if err != nil {
		s.logError(count, totalCount, err)
		return
	}
	defer resp.Body.Close()

	var content []byte
	switch s.matchFrom {
	case "body":
		content, err = io.ReadAll(resp.Body)
	case "header":
		content = []byte(fmt.Sprint(resp.Header))
	}

	if err != nil {
		s.logError(count, totalCount, err)
		return
	}
	s.printResult(url, count, totalCount, resp.StatusCode, content)
}

func (s *Scanner) logError(count uint64, totalCount uint64, err error) {
	fmt.Printf("[ %d / %d ] - [ %s ] - %s\n", count, totalCount, color.Red.AnsiFormat("ERROR"), err)
}

func (s *Scanner) printResult(url string, count, totalCount uint64, statusCode int, content []byte) {
	statusStr := strconv.Itoa(statusCode)
	var statusColor color.Color

	switch {
	case statusCode >= 200 && statusCode < 300:
		statusColor = color.BoldGreen
		if s.m.Match(content) {
			fmt.Printf("[ %d / %d ] - [ %s ] - %s\n", count, totalCount, color.Green.AnsiFormat(s.name), url)
			utils.WriteFile(s.o, []byte(url))
			return
		} else {
			statusColor = color.BoldYellow
		}
	case statusCode >= 300 && statusCode < 400:
		statusColor = color.BoldCyan
	case statusCode >= 400 && statusCode < 500:
		statusColor = color.BoldRed
	case statusCode >= 500:
		statusColor = color.BoldMagenta
	default:
		statusColor = color.White
	}

	fmt.Printf("[ %d / %d ] - [ %s ] - %s\n", count, totalCount, statusColor.AnsiFormat(statusStr), url)
}
