package app

import (
	"fmt"
	"time"
)

// Activity is an entry in the schedule
type Activity struct {
	Type   string
	Hour   int
	Minute int
	Day    time.Weekday
	Daily  bool
}

// Time the Activity starts
func (a *Activity) Time() time.Time {
	now := time.Now()
	dayOffset := 0
	if a.Daily {
		if a.Hour < now.Hour() {
			dayOffset++
		} else if a.Hour == now.Hour() && a.Minute < now.Minute() {
			dayOffset++
		}
	} else {
		dayOffset = int(a.Day - now.Weekday())
		if dayOffset < 0 {
			dayOffset += 7
		} else if a.Hour < now.Hour() {
			dayOffset += 7
		} else if a.Hour == now.Hour() && a.Minute < now.Minute() {
			dayOffset += 7
		}
	}
	return time.Date(now.Year(), now.Month(), now.Day()+dayOffset, a.Hour, a.Minute, 0, 0, now.Location())
}

func (a *Activity) String() string {
	return fmt.Sprintf("%s, %.0f hours", a.Type, time.Until(a.Time()).Hours())

}

// DailyActivity creates a activity
func DailyActivity(t string, hour, minute int) *Activity {
	return &Activity{Daily: true, Type: t, Hour: hour, Minute: minute}
}

// WeeklyActivity creates a activity
func WeeklyActivity(d time.Weekday, t string, hour, minute int) *Activity {
	return &Activity{Type: t, Day: d, Hour: hour, Minute: minute}
}
