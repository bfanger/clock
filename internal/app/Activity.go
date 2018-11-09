package app

import (
	"time"
)

// Activity is an entry in the schedule
type Activity struct {
	Hour   int
	Minute int
	Day    time.Weekday
	Daily  bool
	alarm  Alarm
}

// NextAlarm the Activity starts
func (a *Activity) NextAlarm() *Alarm {
	now := time.Now()
	dayOffset := 0
	if a.Daily {
		if a.Hour < now.Hour() {
			dayOffset++
		} else if a.Hour == now.Hour() && a.Minute <= now.Minute() {
			dayOffset++
		}
	} else {
		dayOffset = int(a.Day - now.Weekday())
		if dayOffset < 0 {
			dayOffset += 7
		} else if a.Hour < now.Hour() {
			dayOffset += 7
		} else if a.Hour == now.Hour() && a.Minute <= now.Minute() {
			dayOffset += 7
		}
	}
	copy := a.alarm
	copy.Start = time.Date(now.Year(), now.Month(), now.Day()+dayOffset, a.Hour, a.Minute, 0, 0, now.Location())
	return &copy
}

// DailyActivity creates a activity
func DailyActivity(hour, minute int, alarm Alarm) *Activity {
	return &Activity{Daily: true, Hour: hour, Minute: minute, alarm: alarm}
}

// WeeklyActivity creates a activity
func WeeklyActivity(d time.Weekday, hour, minute int, alarm Alarm) *Activity {
	return &Activity{Day: d, Hour: hour, Minute: minute, alarm: alarm}
}
