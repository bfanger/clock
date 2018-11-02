// https://inzamelkalender.hvcgroep.nl/ical-info

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bfanger/clock/internal/app"
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
		t, err := nextGarbageTruck()
		if err != nil {
			log.Fatal(err)
		}
		d := time.Until(t.Arrival)
		hours := d.Hours()
		log.Printf("Next %s in %.1f hours\n", t.Type, hours)
		active := hours < hoursBefore

		if first || previous != active {
			previous = active
			first = false
			if !active {
				wait = d - (hoursBefore * time.Hour) + time.Minute
			} else {
				wait = d + (hoursAfter) + time.Minute
			}
			if active {
				log.Printf("Show notification %s for %s\n", t.Type, d+hoursAfter)
				if err := app.ShowNotification(t.Type, d+hoursAfter); err != nil {
					log.Printf("Failed to show notfication: %v", err)
				}
			}
		}

		log.Printf("Sleeping for %.1f hours\n", wait.Hours())
		time.Sleep(wait)
	}
}

type garbageTruck struct {
	Type    string
	Arrival time.Time
}

func nextGarbageTruck() (*garbageTruck, error) {
	events, err := garbageCalendar()
	if err != nil {
		return nil, fmt.Errorf("failed to load events: %v", err)
	}
	t := garbageTruck{}
	for _, e := range events {
		d := time.Until(e.Start)
		if d.Hours() < -1*hoursAfter {
			continue // Skip old entries
		}

		switch e.Summary {
		case "Restafval":
			t.Type = "restafval"
		case "Papier en karton":
			t.Type = "restafval"
		}
		if t.Type != "" {
			t.Arrival = e.Start
			return &t, nil
		}
	}
	return nil, errors.New("no valid entries found")
}

func garbageCalendar() ([]*event, error) {
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
