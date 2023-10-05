package main

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/internal/ical"
	"github.com/bfanger/clock/internal/schedule"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	for {
		appointment, err := nextSchoolDay()
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

func nextSchoolDay() (*schedule.Appointment, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	rawUrl := os.Getenv("SCHOOL_WEBCAL")
	if rawUrl == "" {
		return nil, errors.New("Missing SCHOOL_WEBCAL")
	}
	url, err := url.Parse(os.Getenv("SCHOOL_WEBCAL"))
	if err != nil {
		return nil, err
	}
	url.Scheme = "https"
	r, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	events, err := ical.Parse(r.Body)
	if err != nil {
		return nil, err
	}
	var appointments []*schedule.Appointment

	for _, d := range ical.GroupByDay(events) {
		appointment := &schedule.Appointment{
			Notification: "school",
			At:           d.Date.Add(-50 * time.Minute),
			Duration:     35 * time.Minute,
			Timer:        30 * time.Minute,
		}
		for _, e := range d.Events {
			if strings.Contains(e.Summary, " lo ") {
				appointment.Notification = "gym"
			}
		}
		appointments = append(appointments, appointment)
	}

	planned := schedule.Upcoming(appointments)
	if len(planned) == 0 {
		return nil, errors.New("outdated calender")
	}
	return planned[0], nil

}
