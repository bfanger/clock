// https://www.hvcgroep.nl/zelf-regelen/afvalkalender-maandoverzicht

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/internal/schedule"
	"github.com/pkg/errors"
)

func main() {
	for {
		appointment, err := nextGarbageTruck()
		if err != nil {
			app.Fatal(err)
		}
		appointment.Wait()
		if err := app.ShowAppointment(appointment); err != nil {
			app.Fatal(err)
		}
		time.Sleep(appointment.Duration)
	}
}

func nextGarbageTruck() (*schedule.Appointment, error) {
	events, err := garbageCalendar()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load events")
	}
	var appointments []*schedule.Appointment
	for _, e := range events {

		notification := schedule.FromTill(
			time.Date(e.Start.Year(), e.Start.Month(), e.Start.Day()-1, 10, 0, 0, 0, time.Local),
			time.Date(e.Start.Year(), e.Start.Month(), e.Start.Day(), 9, 0, 0, 0, time.Local),
		)
		switch strings.ToLower(e.Summary) {
		case "plastic\\, blik & drinkpakken":
			notification.Notification = "plastic"
			appointments = append(appointments, notification)
		case "papier en karton":
			notification.Notification = "papier"
			appointments = append(appointments, notification)
		case "gft & etensresten":
		case "gft & etensresten.":
			notification.Notification = "gft"
			appointments = append(appointments, notification)
			continue
		default:
			fmt.Printf("unknown event: %s\n", e.Summary)
		}
	}
	if len(appointments) == 0 {
		return nil, errors.New("no valid entries found")
	}
	planned := schedule.Upcomming(appointments)
	if len(planned) == 0 {
		return nil, errors.New("outdated calender")
	}
	return planned[0], nil
}

type event struct {
	Summary string
	Start   time.Time
}

func garbageCalendar() ([]event, error) {
	const eventMode = "VEVENT"
	r, err := http.Get("https://inzamelkalender.hvcgroep.nl/ical/0479200000012088")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	lines := bufio.NewScanner(r.Body)
	var stack []string
	var mode string
	var events []event
	var e event
	for lines.Scan() {
		line := lines.Text()
		if strings.HasPrefix(line, "BEGIN:") {
			mode = line[6:]
			stack = append(stack, mode)
			if mode == eventMode {
				e = event{}
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
