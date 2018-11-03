package main

import (
	"fmt"
	"time"

	"github.com/bfanger/clock/internal/app"
)

func main() {
	schedule := []*app.Activity{
		app.WeeklyActivity(time.Saturday, "zwemmen", 15, 45),
		app.WeeklyActivity(time.Monday, "school", 8, 05),
		app.WeeklyActivity(time.Tuesday, "gym", 8, 10),
		app.WeeklyActivity(time.Wednesday, "school", 8, 10),
		app.WeeklyActivity(time.Thursday, "gym", 8, 05),
		app.WeeklyActivity(time.Friday, "school", 8, 10),
		app.DailyActivity("vis", 20, 0)}

	for {
		a := nextActivity(schedule)
		t := a.Time()
		fmt.Printf("Scheduled reminder: \"%s\" on %s %d:%02d\n", a.Type, t.Weekday(), t.Hour(), t.Minute())
		time.Sleep(time.Until(t))
		fmt.Printf("Showing notification %s\n", a.Type)
		app.ShowNotification(a.Type, 10*time.Minute)
	}
}

func nextActivity(schedule []*app.Activity) (result *app.Activity) {
	first := time.Now().Add(365 * 24 * time.Hour)
	for _, a := range schedule {
		start := a.Time()
		if start.Before(first) {
			first = start
			result = a
		}
	}
	return
}
