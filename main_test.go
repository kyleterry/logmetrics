package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestParseEntriesCanFormatForStatsdAndWriteToWriter(t *testing.T) {
	entries := []string{
		`192.168.99.1 - - - - - [08/Jul/2016:01:30:19 +0000] http - http "GET / HTTP/1.1" 200 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:30:19 +0000] http - http "GET / HTTP/1.1" 200 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:30:19 +0000] http - http "GET / HTTP/1.1" 200 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:30:20 +0000] http - http "GET / HTTP/1.1" 200 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:30:20 +0000] http - http "GET / HTTP/1.1" 200 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:30:20 +0000] http - http "GET / HTTP/1.1" 200 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:30:20 +0000] http - http "GET / HTTP/1.1" 200 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:30:25 +0000] http - http "GET /error_page HTTP/1.1" 500 2 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:30:25 +0000] http - http "GET /error_page HTTP/1.1" 500 2 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:30:30 +0000] http - http "GET /error_page HTTP/1.1" 500 2 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:31:00 +0000] http - http "GET /error_page HTTP/1.1" 500 2 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:32:00 +0000] http - http "GET /error_page HTTP/1.1" 500 2 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:35:49 +0000] http - http "GET /error_page HTTP/1.1" 500 2 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:35:59 +0000] http - http "GET / HTTP/1.1" 304 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:35:59 +0000] http - http "GET / HTTP/1.1" 304 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:35:59 +0000] http - http "GET / HTTP/1.1" 304 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:35:59 +0000] http - http "GET / HTTP/1.1" 304 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:35:59 +0000] http - http "GET / HTTP/1.1" 304 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:35:59 +0000] http - http "GET / HTTP/1.1" 304 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:36:00 +0000] http - http "GET / HTTP/1.1" 304 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:36:00 +0000] http - http "GET / HTTP/1.1" 304 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:36:00 +0000] http - http "GET / HTTP/1.1" 304 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:36:00 +0000] http - http "GET / HTTP/1.1" 304 0 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:45:57 +0000] http - http "GET /blah HTTP/1.1" 404 209 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:45:58 +0000] http - http "GET /blah HTTP/1.1" 404 209 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:45:58 +0000] http - http "GET /blah HTTP/1.1" 404 209 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:45:58 +0000] http - http "GET /blah HTTP/1.1" 404 209 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:45:58 +0000] http - http "GET /blah HTTP/1.1" 404 209 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:45:58 +0000] http - http "GET /blah HTTP/1.1" 404 209 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:45:59 +0000] http - http "GET /blah HTTP/1.1" 404 209 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:45:59 +0000] http - http "GET /blah HTTP/1.1" 404 209 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:45:59 +0000] http - http "GET /blah HTTP/1.1" 404 209 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
		`192.168.99.1 - - - - - [08/Jul/2016:01:45:59 +0000] http - http "GET /blah HTTP/1.1" 404 209 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"`,
	}

	expected := []string{
		"20x:7|s",
		"30x:10|s",
		"40x:10|s",
		"50x:6|s",
		"/error_page:6|s",
	}

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	parseEntries(entries, w)

	result := buf.String()
	for _, line := range expected {
		if !strings.Contains(result, line) {
			t.Fatal("Expected:", line, "to be in the results")
		}
	}
}
