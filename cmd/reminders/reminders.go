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
			Notification: "school",
			Hour:         7,
			Minute:       45,
			Duration:     45 * time.Minute,
			Timer:        30 * time.Minute,
			Repeat:       schedule.RepeatDays{Tuesday: true, Wednesday: true, Friday: true},
		},
		{
			Notification: "gym",
			Hour:         7,
			Minute:       45,
			Duration:     45 * time.Minute,
			Timer:        30 * time.Minute,
			Repeat:       schedule.RepeatDays{Monday: true, Thursday: true},
		},
		{
			Notification: "zwemmen",
			Hour:         15,
			Minute:       35,
			Duration:     15 * time.Minute,
			Timer:        15 * time.Minute,
			Repeat:       schedule.RepeatDays{Saturday: true, Thursday: true},
		},
		{
			Notification: "bedtime-charlie",
			Hour:         20,
			Minute:       10,
			Duration:     15 * time.Minute,
			Repeat:       schedule.Daily(),
		},
		{
			Notification: "bedtime-bob",
			Hour:         23,
			Minute:       45,
			Duration:     45 * time.Minute,
			Repeat:       schedule.Daily(),
		},
	}
	for {
		appointments := schedule.PlanRepeated(schema)
		if len(appointments) == 0 {
			app.Fatal(errors.New("Empty schedule"))
		}
		planned := schedule.Upcomming(appointments)
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
