package ical

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sort"
	"time"
)

type Event struct {
	Summary     string
	Description string
	Location    string
	Uid         string
	Start       time.Time
	End         time.Time
}

func Parse(r io.ReadCloser) ([]Event, error) {
	defer r.Close()
	s := &Scanner{source: bufio.NewScanner(r)}
	var events []Event
	var event *Event
	for s.Scan() {
		if s.Key == "BEGIN" && s.Value == "VEVENT" {
			event = &Event{}
			continue
		}
		if s.Key == "END" && s.Value == "VEVENT" {
			events = append(events, *event)
			event = nil
			continue
		}
		if event == nil {
			continue
		}
		switch s.Key {
		case "SUMMARY":
			event.Summary = s.Value
		case "DESCRIPTION":
			event.Description = s.Value
		case "LOCATION":
			event.Location = s.Value
		case "UID":
			event.Uid = s.Value
		case "DTSTART":
			start, err := time.Parse("20060102T150405Z", s.Value)
			if err != nil {
				return nil, err
			}
			event.Start = start
		case "DTEND":
			end, err := time.Parse("20060102T150405Z", s.Value)
			if err != nil {
				return nil, err
			}
			event.End = end
		}
	}
	if s.Err != nil {
		return nil, s.Err
	}

	sort.SliceStable(events, func(i, j int) bool { return events[i].Start.Before(events[j].Start) })
	return events, nil
}

type Scanner struct {
	Key    string
	Value  string
	source *bufio.Scanner
	buffer string
	Err    error
}

var r = regexp.MustCompile(`^([^:]+):(.*)$`)

func (s *Scanner) Scan() bool {
	if s.Err != nil {
		return false
	}
	line := s.buffer
	more := true

	for {
		if !s.source.Scan() {
			s.Err = s.source.Err()
			more = false
			break
		}
		next := s.source.Text()
		if line == "" {
			line = next
			continue
		}
		if next[0] == ' ' {
			// multiline value
			line += next[1:]
			continue
		}
		s.buffer = next
		break
	}
	matches := r.FindStringSubmatch(line)
	if len(matches) != 3 {
		s.Err = fmt.Errorf("invalid line: %s", line)
		return false
	}
	s.Key = matches[1]
	s.Value = matches[2]

	return more
}

type PerDay struct {
	Date   time.Time
	Events []Event
}

func GroupByDay(events []Event) []PerDay {
	var grouped []PerDay
	for _, c := range events {
		var day *PerDay
		for _, d := range grouped {
			if d.Date.Format("20060102") == c.Start.Format("20060102") {
				day = &d
			}
		}
		if day == nil {
			grouped = append(grouped, PerDay{
				Date:   c.Start,
				Events: []Event{c},
			})
		} else {
			day.Events = append(day.Events, c)
		}
	}
	return grouped
}
