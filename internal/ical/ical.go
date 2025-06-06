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

func Parse(r io.Reader) ([]Event, error) {
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
			start, err := s.Date()
			if err != nil {
				return nil, err
			}
			event.Start = start
		case "DTEND":
			end, err := s.Date()
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
	Meta   string
	Value  string
	source *bufio.Scanner
	buffer string
	Err    error
}

var (
	keyValueSplit = regexp.MustCompile(`^([^:]+):(.*)$`)
	keyMetaSplit  = regexp.MustCompile(`^([^;]+);(.*)$`)
)

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
	matches := keyValueSplit.FindStringSubmatch(line)
	if len(matches) != 3 {
		s.Err = fmt.Errorf("invalid line: %s", line)
		return false
	}
	s.Key = matches[1]
	s.Value = matches[2]
	matches = keyMetaSplit.FindStringSubmatch(s.Key)
	if len(matches) == 3 {
		s.Key = matches[1]
		s.Meta = matches[2]
	}
	return more
}

func (s *Scanner) Date() (time.Time, error) {
	if len(s.Value) == 8 {
		return time.Parse("20060102", s.Value)
	}
	return time.Parse("20060102T150405Z", s.Value)
}

type PerDay struct {
	Date   time.Time
	Events []Event
}

func GroupByDay(events []Event) []PerDay {
	var grouped []PerDay

	for _, event := range events {
		var day *PerDay
		for i, d := range grouped {
			if d.Date.Local().Format(time.DateOnly) == event.Start.Local().Format(time.DateOnly) {
				day = &grouped[i]
				break
			}
		}
		if day == nil {
			grouped = append(grouped, PerDay{
				Date:   event.Start,
				Events: []Event{event},
			})
		} else {
			day.Events = append(day.Events, event)
		}
	}
	return grouped
}
