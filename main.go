package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var (
	inputFile  string
	outputFile string
	interval   int
)

func init() {
	flag.StringVar(&inputFile, "input-file", "/var/log/nginx/access.log", "Nginx log file location")
	flag.StringVar(&outputFile, "output-file", "/var/log/stats.log", "Output log file location")
	flag.IntVar(&interval, "interval", 5, "Time (in seconds) to scan for new log entries")
}

func periodicScan(file *os.File, out io.Writer, quit chan struct{}) {
	// seek to the end of the file before we start scanning
	file.Seek(0, os.SEEK_END)

	for {
		select {
		case <-time.After(time.Second * time.Duration(interval)):
			scanner := bufio.NewScanner(file)
			entries := []string{}
			for scanner.Scan() {
				entries = append(entries, scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			if len(entries) > 0 { // no reason to parse if we don't have entries
				go parseEntries(entries, out)
			}
		case <-quit:
			break
		}
	}
}

func parseEntries(entries []string, out io.Writer) {
	stats := make(map[string]int)
	w := bufio.NewWriter(out)
	for _, entry := range entries {
		parts := strings.Fields(entry)
		statusCode, err := strconv.Atoi(parts[14])
		if err != nil {
			log.Println(err)
			continue
		}
		route := parts[12]
		switch {
		case statusCode > 199 && statusCode < 300:
			stats["20x"] += 1
		case statusCode > 299 && statusCode < 400:
			stats["30x"] += 1
		case statusCode > 399 && statusCode < 500:
			stats["40x"] += 1
		case statusCode > 499 && statusCode < 600:
			stats["50x"] += 1
			stats[route] += 1
		default:
			continue
		}
	}
	for k, v := range stats {
		w.WriteString(fmt.Sprintf("%s:%d|s\n", k, v))
	}
	w.Flush()
}

func main() {
	flag.Parse()

	log.Printf("Starting up... Scan interval is set to %d seconds", interval)

	in, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	quit := make(chan struct{})
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt)

	go func() {
		<-sch
		close(quit)
	}()

	go periodicScan(in, out, quit)

	<-quit
}
