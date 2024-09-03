package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"go-exposed-config-scanner/pkg/color"
	"go-exposed-config-scanner/pkg/matcher"
	"go-exposed-config-scanner/pkg/request"
	"go-exposed-config-scanner/pkg/utils"
)

func NewScanner(c *request.Requester, m matcher.Matcher, o *os.File, n string, mf string, v bool, mo bool) *Scanner {
	return &Scanner{
		client:    c,
		matcher:   m,
		output:    o,
		name:      n,
		matchFrom: mf,
		verbose:   v,
		matchOnly: mo,
	}
}

func (s *Scanner) getContent(r *http.Response) ([]byte, error) {
	switch strings.ToLower(s.matchFrom) {
	case "body":
		return io.ReadAll(r.Body)
	case "headers":
		return []byte(fmt.Sprint(r.Header)), nil
	}
	return nil, fmt.Errorf("invalid matchFrom value: %s", s.matchFrom)
}

func (s *Scanner) Scan(url string, count uint64, totalCount uint64) {
	resp, err := s.client.Do(url)
	if err != nil {
		s.logError(count, totalCount, url, err)
		return
	}
	defer resp.Body.Close()

	content, err := s.getContent(resp)
	if err != nil {
		s.logError(count, totalCount, url, err)
		return
	}

	matched := s.matcher.Match(content)
	s.logResult(count, totalCount, url, matched)

	if matched {
		utils.WriteFile(s.output, []byte(url))
	}
}

func (s *Scanner) logResult(count uint64, totalCount uint64, url string, matched bool) {
	if s.matchOnly && !matched {
		return
	}

	status := color.Red.AnsiFormat("BAD")
	if matched {
		status = color.Green.AnsiFormat("OK")
	}

	output := fmt.Sprintf("[%s/%s] %s %s %s - %s",
		color.Cyan.AnsiFormat(fmt.Sprintf("%d", count)),
		color.Cyan.AnsiFormat(fmt.Sprintf("%d", totalCount)),
		color.Yellow.AnsiFormat(s.name),
		color.White.AnsiFormat("-"),
		color.Blue.AnsiFormat(url),
		status,
	)

	fmt.Println(output)
}

func (s *Scanner) logError(count uint64, totalCount uint64, url string, err error) {
	if !s.verbose {
		return
	}

	errorOutput := fmt.Sprintf("[%s/%s] %s %s %s - %s",
		color.Cyan.AnsiFormat(fmt.Sprintf("%d", count)),
		color.Cyan.AnsiFormat(fmt.Sprintf("%d", totalCount)),
		color.Yellow.AnsiFormat(s.name),
		color.White.AnsiFormat("-"),
		color.Blue.AnsiFormat(url),
		color.Red.AnsiFormat(err.Error()))

	fmt.Println(errorOutput)
}
