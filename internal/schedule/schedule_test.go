package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRepeating(t *testing.T) {
	datetime := time.Date(2010, 10, 5, 14, 30, 0, 0, time.Local)
	a := daily(16, 30)
	assert.Equal(t, a.plannedAfter(datetime), datetime.Add(2*time.Hour), "appointment start in 2 hours")
	assert.Equal(t, a.plannedAfter(datetime.Add(3*time.Hour)), datetime.Add(26*time.Hour), "appointment today was 1 hour ago, plan the next day")
	a.Repeat.Tuesday = false
	assert.Equal(t, a.plannedAfter(datetime), datetime.Add(26*time.Hour), "not scheduled today, plan the next day")
}

func daily(hour, minute int) Appointment {
	a := Appointment{
		Hour:   hour,
		Minute: minute,
	}
	a.Repeat.Monday = true
	a.Repeat.Tuesday = true
	a.Repeat.Wednesday = true
	a.Repeat.Thursday = true
	a.Repeat.Friday = true
	a.Repeat.Saturday = true
	a.Repeat.Sunday = true
	return a
}
