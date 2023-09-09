package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRepeating(t *testing.T) {
	datetime := time.Date(2010, 10, 5, 14, 30, 0, 0, time.Local)
	appointment := &RepeatedAppointment{Hour: 16, Minute: 30, Repeat: Daily()}
	actual, err := appointment.plannedAfter(datetime)
	assert.Nil(t, err)
	assert.Equal(t, actual.At, datetime.Add(2*time.Hour), "appointment starts in 2 hours")
	actual, err = appointment.plannedAfter(datetime.Add(3 * time.Hour))
	assert.Nil(t, err)
	assert.Equal(t, actual.At, datetime.Add(26*time.Hour), "appointment today was 1 hour ago, planned the next day")
	appointment.Repeat.Tuesday = false
	actual, err = appointment.plannedAfter(datetime)
	assert.Nil(t, err)
	assert.Equal(t, actual.At, datetime.Add(26*time.Hour), "not scheduled today, planned the next day")
}
func TestSchedule(t *testing.T) {
	schema := []*RepeatedAppointment{
		{Notification: "half vijf", Hour: 16, Minute: 30, Repeat: Daily()},
		{Notification: "acht uur", Hour: 20, Minute: 0, Repeat: Daily()},
	}

	datetime := time.Date(2010, 10, 5, 14, 30, 0, 0, time.Local)
	appointments := planRepeatedAfter(schema, datetime)
	assert.Len(t, appointments, 2)
	planned := upcomingAfter(appointments, datetime)
	assert.Len(t, planned, 1)
	assert.Equal(t, "half vijf", planned[0].Notification)
	planned = upcomingAfter(appointments, datetime.Add(3*time.Hour))
	assert.Len(t, planned, 1)
	assert.Equal(t, "acht uur", planned[0].Notification)
	schema = append(schema, &RepeatedAppointment{Notification: "duplicate time", Hour: 16, Minute: 30, Repeat: Daily()})
	appointments = planRepeatedAfter(schema, datetime)
	assert.Len(t, appointments, 3)
	planned = upcomingAfter(appointments, datetime)
	assert.Len(t, planned, 2)
}
