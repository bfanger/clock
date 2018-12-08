package app

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Alarm is the notification trigger
type Alarm struct {
	Notification string
	Start        time.Time
	Duration     time.Duration
}

const endpoint = "http://localhost:8080/"

// Activate the alarm
func (a *Alarm) Activate() error {
	fmt.Printf("Showing notification %s\n", a.Notification)
	data := url.Values{}
	data.Set("action", "notify")
	data.Set("icon", a.Notification)
	data.Set("duration", strconv.Itoa(int(a.Duration.Seconds())))
	r, err := http.PostForm(endpoint+"notify", data)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

func (a *Alarm) String() string {
	return a.Notification + " notification on " + a.Start.Format("Mon 2 January 15:04")
}

// FirstAlarm returns the earliest alarm in the future
func FirstAlarm(alarms []*Alarm) *Alarm {
	var first *Alarm
	now := time.Now()
	for _, a := range alarms {
		if a.Start.Before(now) {
			continue
		}
		if first == nil || a.Start.Before(first.Start) {
			first = a
		}
	}
	return first
}
