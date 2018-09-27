// https://inzamelkalender.hvcgroep.nl/ical-info

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type event struct {
	Summary string
	Start   time.Time
}

func main() {
	events, err := getEvents()
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now().Add(-2 * time.Hour)
	end := start.Add(24 * time.Hour)
	for _, e := range events {
		if e.Start.After(start) && e.Start.Before(end) {
			fmt.Printf("%s: %s\n", e.Summary, e.Start)
			res, err := http.PostForm("http://localhost:8080/", url.Values{"action": {"Show"}})
			if err != nil {
				fmt.Println(err)
				continue
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(string(body))
		}
	}
}

func getEvents() ([]*event, error) {
	const eventMode = "VEVENT"

	r, err := http.Get("https://inzamelkalender.hvcgroep.nl/ical/0479200000012088")
	if err != nil {
		return nil, err
	}
	lines := bufio.NewScanner(r.Body)
	var stack []string
	var mode string
	var events []*event
	var e *event

	for lines.Scan() {
		line := lines.Text()
		if strings.HasPrefix(line, "BEGIN:") {
			mode = line[6:]
			stack = append(stack, mode)
			if mode == eventMode {
				e = &event{}
			}
		}
		if strings.HasPrefix(line, "END:") {
			stack = stack[:len(stack)-1]
			mode = stack[len(stack)-1]

			if mode == eventMode {
				events = append(events, e)
			}
		}
		if mode == eventMode {
			if strings.HasPrefix(line, "SUMMARY:") {
				e.Summary = line[8:]
			}
			if strings.HasPrefix(line, "DTSTART:") {
				e.Start, err = time.Parse("20060102T150405", line[8:])
			}
		}
	}
	if err := lines.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
