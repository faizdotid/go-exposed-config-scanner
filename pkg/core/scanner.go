package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"go-exposed-config-scanner/pkg/color"
	"go-exposed-config-scanner/pkg/matcher"
	"go-exposed-config-scanner/pkg/request"
	"go-exposed-config-scanner/pkg/utils"
)

func NewScanner(c *request.Requester, m matcher.Matcher, o *os.File, n string, mf string, v bool, mo bool, counter *atomic.Uint64, totalCount uint64, mu *sync.Mutex) *Scanner {
	return &Scanner{
		client:     c,
		matcher:    m,
		output:     o,
		name:       n,
		matchFrom:  mf,
		verbose:    v,
		matchOnly:  mo,
		counter:    counter,
		totalCount: totalCount,
		mu:         mu,
	}
}

func (s *Scanner) matchContent(r *http.Response) (bool, error) {
	switch strings.ToLower(s.matchFrom) {
	case "body":
		content, err := io.ReadAll(r.Body)
		if err != nil {
			return false, err
		}
		return s.matcher.Match(content), nil
	case "headers":
		for _, value := range r.Header {
			if s.matcher.Match([]byte(strings.Join(value, ","))) {
				return true, nil
			}
		}
	default:
		return false, fmt.Errorf("invalid matchFrom value: %s", s.matchFrom)
	}
	return false, nil
}

func (s *Scanner) Scan(url string) {
	resp, err := s.client.Do(url)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter.Add(1)
	if err != nil {
		s.logError(url, err)
		return
	}
	defer resp.Body.Close()

	matched, err := s.matchContent(resp)
	if err != nil {
		s.logError(url, err)
		return
	}
	s.logResult(url, matched)

	if matched {
		utils.WriteFile(s.output, []byte(url))
	}
}

func (s *Scanner) logResult(url string, matched bool) {
	if s.matchOnly && !matched {
		return
	}

	status := color.Red.AnsiFormat("BAD")
	if matched {
		status = color.Green.AnsiFormat("OK")
	}

	output := fmt.Sprintf("[%s/%s] %s %s %s - %s",
		color.Cyan.AnsiFormat(fmt.Sprintf("%d", s.counter.Load())),
		color.Cyan.AnsiFormat(fmt.Sprintf("%d", s.totalCount)),
		color.Yellow.AnsiFormat(s.name),
		color.White.AnsiFormat("-"),
		color.Blue.AnsiFormat(url),
		status,
	)

	fmt.Println(output)
}

func (s *Scanner) logError(url string, err error) {
	if !s.verbose {
		return
	}

	errorOutput := fmt.Sprintf("[%s/%s] %s %s %s - %s",
		color.Cyan.AnsiFormat(fmt.Sprintf("%d", s.counter.Load())),
		color.Cyan.AnsiFormat(fmt.Sprintf("%d", s.totalCount)),
		color.Yellow.AnsiFormat(s.name),
		color.White.AnsiFormat("-"),
		color.Blue.AnsiFormat(url),
		color.Red.AnsiFormat(err.Error()))

	fmt.Println(errorOutput)
}
