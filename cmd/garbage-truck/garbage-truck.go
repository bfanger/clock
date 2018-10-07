// https://inzamelkalender.hvcgroep.nl/ical-info

package main

import (
	"bufio"
	"errors"
	"fmt"
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

const hoursBefore = 13 // Start notification at 20:00 the day before
const hoursAfter = 2   // Hide the notification at 9:00

func main() {
	first := true
	var previous bool
	var wait time.Duration
	for {
		d, err := nextGarbageTruck()
		if err != nil {
			log.Fatal(err)
		}
		hours := d.Hours()
		log.Printf("Next garbage truck in %.1f hours\n", hours)
		active := hours < hoursBefore

		if first || previous != active {
			previous = active
			first = false

			if active {
				log.Println("Show notification")
			} else {
				log.Println("Hide notification")
			}
			if err := notify(active); err != nil {
				log.Printf("Failed to send notfication: %v", err)
			}

		}
		if !active {
			wait = d - (hoursBefore * time.Hour) + time.Minute
		} else {
			wait = d + (hoursAfter) + time.Minute
		}
		log.Printf("Sleeping for %.1f hours\n", wait.Hours())
		time.Sleep(wait)
	}

}
func nextGarbageTruck() (time.Duration, error) {
	events, err := garbageCalender()
	if err != nil {
		return 0, fmt.Errorf("failed to load events: %v", err)
	}
	for _, e := range events {
		if e.Summary != "Restafval" {
			continue
		}
		d := time.Until(e.Start)
		if d.Hours() < -1*hoursAfter {
			continue // Skip old entries
		}
		return d, nil
	}
	return 0, errors.New("No Restafval entries found")
}

func garbageCalender() ([]*event, error) {
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
			if len(stack) == 0 {
				mode = ""
			} else {
				mode = stack[len(stack)-1]
			}

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

func notify(toggle bool) error {
	data := url.Values{}
	if toggle {
		data.Set("action", "Show")
	} else {
		data.Set("action", "Hide")
	}
	if _, err := http.PostForm("http://localhost:8080/", data); err != nil {
		return err
	}
	return nil
}
