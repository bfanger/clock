package main

import (
	"fmt"
	"time"

	"github.com/bfanger/clock/internal/app"
)

func main() {
	schedule := []*app.Activity{
		app.WeeklyActivity(time.Saturday, 15, 45, app.Alarm{Notification: "zwemmen", Duration: 15 * time.Minute}),
		app.WeeklyActivity(time.Monday, 7, 35, app.Alarm{Notification: "school", Duration: 45 * time.Minute}),
		app.WeeklyActivity(time.Tuesday, 7, 45, app.Alarm{Notification: "gym", Duration: 45 * time.Minute}),
		app.WeeklyActivity(time.Wednesday, 7, 45, app.Alarm{Notification: "school", Duration: 45 * time.Minute}),
		app.WeeklyActivity(time.Thursday, 7, 35, app.Alarm{Notification: "gym", Duration: 45 * time.Minute}),
		app.WeeklyActivity(time.Friday, 7, 45, app.Alarm{Notification: "school", Duration: 45 * time.Minute}),
		app.DailyActivity(20, 0, app.Alarm{Notification: "vis", Duration: 10 * time.Minute}),
		app.DailyActivity(0, 0, app.Alarm{Notification: "bedtime", Duration: 30 * time.Minute}),
	}

	for {
		alarms := make([]*app.Alarm, len(schedule))
		for i, activity := range schedule {
			alarms[i] = activity.NextAlarm()
		}
		alarm := app.FirstAlarm(alarms)
		fmt.Printf("Next: %s\n", alarm)
		time.Sleep(time.Until(alarm.Start))
		alarm.Activate()
		time.Sleep(alarm.Duration)
	}
}
