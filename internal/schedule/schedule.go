package schedule

import (
	"log"
	"time"

	"github.com/pkg/errors"
)

// Appointment is a single occurance
type Appointment struct {
	Notification string
	At           time.Time
	Duration     time.Duration
	Timer        time.Duration
}

// Wait until the appointment
func (a *Appointment) Wait() {
	wait := time.Until(a.At)
	log.Printf("%s", a.At)
	log.Printf("wait: '%s' on %s", a.Notification, a.At.Format("Mon 2 January 15:04"))
	time.Sleep(wait)
}

// RepeatedAppointment of an event, which can trigger a notification
type RepeatedAppointment struct {
	Notification string
	Hour         int // hour
	Minute       int // minutes
	Duration     time.Duration
	Timer        time.Duration
	Repeat       RepeatDays
}

// RepeatDays configuration for an appointment
type RepeatDays struct {
	Monday    bool
	Tuesday   bool
	Wednesday bool
	Thursday  bool
	Friday    bool
	Saturday  bool
	Sunday    bool
}

// Daily repeats every day
func Daily() RepeatDays {
	return RepeatDays{
		Monday:    true,
		Tuesday:   true,
		Wednesday: true,
		Thursday:  true,
		Friday:    true,
		Saturday:  true,
		Sunday:    true,
	}
}

// Planned calculates first occurrence of the appointment
func (a *RepeatedAppointment) Planned() (*Appointment, error) {
	return a.plannedAfter(time.Now())
}

// Planned calculates first occurrence of the appointment after the given time
func (a *RepeatedAppointment) plannedAfter(after time.Time) (*Appointment, error) {
	day := after.Day()
	weekdays := a.repeatedWeekdays()
	for i := 0; i <= 7; i++ {
		planned := time.Date(after.Year(), after.Month(), day+i, a.Hour, a.Minute, 0, 0, time.Local)
		if weekdays[planned.Weekday()] {
			if planned.After(after) {
				return &Appointment{
					Notification: a.Notification,
					At:           planned,
					Duration:     a.Duration,
					Timer:        a.Timer,
				}, nil
			}
		}
	}
	return nil, errors.New("appointment is not repeated on any day")
}

func (a *RepeatedAppointment) repeatedWeekdays() map[time.Weekday]bool {
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

// PlanRepeated converts a repeatable appointments
func PlanRepeated(schema []*RepeatedAppointment) []*Appointment {
	return planRepeatedAfter(schema, time.Now())
}
func planRepeatedAfter(schema []*RepeatedAppointment, after time.Time) []*Appointment {
	var planned []*Appointment
	for _, a := range schema {
		appointment, err := a.plannedAfter(after)
		if err != nil {
			// skip appointment if it couldn't be planned
			continue
		}
		planned = append(planned, appointment)
	}
	return planned
}

// Upcomming appointment(s), returns multiple if they are starting at the same time
func Upcomming(appointments []*Appointment) []*Appointment {
	return upcommingAfter(appointments, time.Now())
}
func upcommingAfter(appointments []*Appointment, datetime time.Time) []*Appointment {
	var upcomming []*Appointment
	for _, a := range appointments {
		if a.At.Before(datetime) {
			continue
		}
		if len(upcomming) == 0 || a.At.Before(upcomming[0].At) {
			upcomming = []*Appointment{a}
		} else if upcomming[0].At == a.At {
			upcomming = append(upcomming, a)
		}
	}
	return upcomming
}
