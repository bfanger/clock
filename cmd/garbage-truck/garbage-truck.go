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

func main() {
	for {
		alarm, err := nextGarbageTruck()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Next: %s\n", alarm)
		time.Sleep(time.Until(alarm.Start))
		alarm.Activate()
		time.Sleep(alarm.Duration)
	}
}

func nextGarbageTruck() (*app.Alarm, error) {
	const hoursBefore = 5 * time.Hour // Start notification at 19:00 the day before
	const hoursAfter = 9 * time.Hour  // Hide the notification at 9:00

	events, err := garbageCalendar()
	if err != nil {
		return nil, fmt.Errorf("failed to load events: %v", err)
	}
	var alarms []*app.Alarm
	for _, e := range events {
		if e.Summary == "Gft & etensresten." {
			continue
		}
		alarm := app.Alarm{
			Start:    e.Start.Add(-1 * hoursBefore),
			Duration: hoursBefore + hoursAfter}

		switch e.Summary {
		case "Restafval":
			alarm.Notification = "restafval"
			alarms = append(alarms, &alarm)
		case "Papier en karton":
			alarm.Notification = "papier"
			alarms = append(alarms, &alarm)
		default:
			fmt.Printf("Unknown event: %s", e.Summary)
		}
	}
	if len(alarms) == 0 {
		return nil, errors.New("no valid entries found")
	}
	return app.FirstAlarm(alarms), nil
}

type event struct {
	Summary string
	Start   time.Time
}

func garbageCalendar() ([]*event, error) {
	const eventMode = "VEVENT"
	r, err := http.Get("https://inzamelkalender.hvcgroep.nl/ical/0479200000012088")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
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
			if strings.HasPrefix(line, "DTSTART;VALUE=DATE:") {
				e.Start, err = time.Parse("20060102", line[19:])
			}
		}
	}
	if err := lines.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
