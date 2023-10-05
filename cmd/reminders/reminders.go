package main

import (
	"errors"
	"time"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/internal/schedule"
)

func main() {
	schema := []*schedule.RepeatedAppointment{
		{
			Notification: "douche",
			Hour:         19,
			Minute:       30,
			Duration:     30 * time.Minute,
			Repeat:       schedule.RepeatDays{Monday: true, Wednesday: true, Friday: true},
		},
		{
			Notification: "bedtime-charlie",
			Hour:         20,
			Minute:       15,
			Duration:     30 * time.Minute,
			Repeat:       schedule.Daily(),
		},
		{
			Notification: "bedtime-bob",
			Hour:         23,
			Minute:       45,
			Duration:     45 * time.Minute,
			Repeat:       schedule.Daily(),
		},
		{
			Notification: "sport",
			Hour:         9,
			Minute:       30,
			Duration:     45 * time.Minute,
			Repeat:       schedule.RepeatDays{Saturday: true},
		},
	}
	for {
		appointments := schedule.PlanRepeated(schema)
		if len(appointments) == 0 {
			app.Fatal(errors.New("Empty schedule"))
		}
		planned := schedule.Upcoming(appointments)
		if len(planned) == 0 {
			app.Fatal(errors.New("No appointments left"))
		}
		planned[0].Wait()
		for _, a := range planned {
			if err := app.ShowAppointment(a); err != nil {
				app.Fatal(err)
			}
		}
	}
}
