package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"go-exposed-config-scanner/pkg/color"
	"go-exposed-config-scanner/pkg/matcher"
	"go-exposed-config-scanner/pkg/request"
	"go-exposed-config-scanner/pkg/utils"
)

func NewScanner(c *request.Requester, m matcher.IMatcher, o string, n string, mf string, v bool, mo bool, counter *atomic.Uint64, totalCount uint64) *Scanner {
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
	}
}

func (s *Scanner) Scan(ctx context.Context, url string, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	resp, err := s.client.Do(ctx, url)
	s.counter.Add(1)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.logError(url, fmt.Errorf("timeout"))
			return
		}

		s.logError(url, err)
		return
	}
	defer resp.Body.Close()

	matched, err := s.matcher.Match(resp)
	if err != nil {
		s.logError(url, err)
		return
	}
	s.logResult(url, matched)

	if matched {
		utils.WriteResultToFile(s.output, url)
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
		color.Red.AnsiFormat(err.Error()),
	)

	fmt.Println(errorOutput)
}
