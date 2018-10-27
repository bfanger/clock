package main

import (
	"log"
	"time"

	"github.com/bfanger/clock/internal/app"
)

func main() {
	schedule := []*app.Activity{
		// WeeklyActivity(time.Saturday, "zwemmen", 15, 50),
		// WeeklyActivity(time.Monday, "school", 8, 15)),
		app.DailyActivity("vis", 20, 0)}

	// schedule = append(schedule, app.DailyActivity("vis", time.Now().Hour(), time.Now().Minute()+1))

	// app.ShowNotification("vis")
	// NewActivity("gym", time.Tuesday, 8, 10)}
	for {
		a := nextActivity(schedule)
		d := time.Until(a.Time())
		log.Printf("Sleeping %s\n", d)
		time.Sleep(d)
		log.Printf("Showing notification %s\n", a.Type)
		app.ShowNotification(a.Type)
		time.Sleep(10 * time.Minute)
		log.Println("Hiding notification")
		app.HideNotification()
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
