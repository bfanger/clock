// https://www.hvcgroep.nl/zelf-regelen/afvalkalender-maandoverzicht

package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/internal/ical"
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
	r, err := http.Get("https://inzamelkalender.hvcgroep.nl/ical/0479200000012088")
	if err != nil {
		return nil, errors.Wrap(err, "failed to load events")
	}
	events, err := ical.Parse(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse calendar")
	}
	var appointments []*schedule.Appointment
	for _, e := range events {

		notification := schedule.FromTill(
			time.Date(e.Start.Year(), e.Start.Month(), e.Start.Day()-1, 12, 0, 0, 0, time.Local),
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
	planned := schedule.Upcoming(appointments)
	if len(planned) == 0 {
		return nil, errors.New("outdated calender")
	}
	return planned[0], nil
}
