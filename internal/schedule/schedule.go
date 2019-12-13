package schedule

import (
	"time"
)

// Appointment of an event, which can trigger a notification
type Appointment struct {
	Disabled    bool `json:"Disabled,omitempty"`
	Description string
	Visual      string
	Hour        int // hour
	Minute      int // minutes
	Duration    int `json:"Duration,omitempty"` // in seconds
	Timer       int `json:"Timer,omitempty"`    // in seconds
	Repeat      struct {
		Monday    bool
		Tuesday   bool
		Wednesday bool
		Thursday  bool
		Friday    bool
		Saturday  bool
		Sunday    bool
	}
}

// Planned calculates first occurrence of the appointment
func (a *Appointment) Planned() time.Time {
	return a.plannedAfter(time.Now())
}

// Planned calculates first occurrence of the appointment after the given time
func (a *Appointment) plannedAfter(after time.Time) time.Time {
	if a.Disabled {
		return time.Time{} //error? appointment was disabled
	}
	daytime := time.Duration(a.Hour)*time.Hour + time.Duration(a.Minute)*time.Minute
	day := after.Day()
	daytimeafter := time.Duration(after.Hour())*time.Hour + time.Duration(after.Minute())*time.Minute
	if daytime < daytimeafter { // not today?
		day++ // check tomorrow
	}
	weekdays := a.repeatedWeekdays()
	for i := 0; i < 7; i++ {
		planned := time.Date(after.Year(), after.Month(), day+i, a.Hour, a.Minute, 0, 0, time.Local)
		if weekdays[planned.Weekday()] {
			return planned
		}
	}
	return time.Time{} // error? appointment is not repeated on any weekday
}

func (a *Appointment) repeatedWeekdays() map[time.Weekday]bool {
	return map[time.Weekday]bool{
		time.Monday:    a.Repeat.Monday,
		time.Tuesday:   a.Repeat.Tuesday,
		time.Wednesday: a.Repeat.Wednesday,
		time.Thursday:  a.Repeat.Thursday,
		time.Friday:    a.Repeat.Friday,
		time.Saturday:  a.Repeat.Saturday,
		time.Sunday:    a.Repeat.Sunday,
	}
}

// Scheduler detemines the next appoinment and prevents duplicates
type Scheduler struct {
	Appointments []*Appointment
	Current      PlannedAppointment
}

type PlannedAppointment struct {
	Appointment *Appointment
	Planned     time.Time
}

func (s *Scheduler) Next() PlannedAppointment {
	var first PlannedAppointment

	now := time.Now()
	for _, a := range s.Appointments {
		p := a.Planned()
		if p.Before(now) {
			continue
		}
		if s.Current.Appointment == a && p == s.Current.Planned {
			p = a.plannedAfter(p.Add(1 * time.Minute))
		}
		if first.Appointment == nil || p.Before(first.Planned) {
			first.Appointment = a
			first.Planned = p
		}
	}
	s.Current = first
	return first
}
