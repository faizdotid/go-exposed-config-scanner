package core

import (
	"fmt"
	"go-exposed-config-scanner/pkg/utils"
	"go-exposed-config-scanner/pkg/validators"
	"io"
	"net/http"
	"os"
)

type Scanner struct {
	c         *http.Client
	v         validators.ValidatorFunction
	o         *os.File
	name      string
	matchFrom string
}

func NewScanner(client *http.Client, validator validators.ValidatorFunction, output *os.File, name string, matchFrom string) *Scanner {
	return &Scanner{
		c:         client,
		v:         validator,
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

	valid, err := s.v(content)
	if err != nil {
		s.logError(count, totalCount, err)
		return
	}

	if valid {
		utils.WriteFile(s.o, []byte(url))
		s.logResult(count, totalCount, "OK", url)
	} else {
		s.logResult(count, totalCount, "FAIL", url)
	}
}

func (s *Scanner) logError(count uint64, totalCount uint64, err error) {
	fmt.Printf("[ %d / %d ] %s - ERR - %s\n", count, totalCount, s.name, err.Error())
}

func (s *Scanner) logResult(count uint64, totalCount uint64, status, url string) {
	fmt.Printf("[ %d / %d ] %s - %s - %s\n", count, totalCount, s.name, status, url)
}