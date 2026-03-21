package timeutil

import (
	"fmt"
	"strings"
	"time"
)

const DefaultLookback = time.Hour

type Range struct {
	From time.Time
	To   time.Time
}

func ParseRange(last, from, to string, now func() time.Time) (Range, error) {
	if now == nil {
		now = time.Now
	}
	if last != "" && (from != "" || to != "") {
		return Range{}, fmt.Errorf("--last cannot be used with --from or --to")
	}

	current := now().UTC()
	if last != "" {
		duration, err := time.ParseDuration(last)
		if err != nil {
			return Range{}, fmt.Errorf("parse --last: %w", err)
		}
		if duration <= 0 {
			return Range{}, fmt.Errorf("--last must be greater than 0")
		}
		return Range{From: current.Add(-duration), To: current}, nil
	}

	toTime, err := parseMoment(to, current)
	if err != nil {
		return Range{}, fmt.Errorf("parse --to: %w", err)
	}
	if toTime.IsZero() {
		toTime = current
	}

	fromTime, err := parseMoment(from, current)
	if err != nil {
		return Range{}, fmt.Errorf("parse --from: %w", err)
	}
	if fromTime.IsZero() {
		fromTime = toTime.Add(-DefaultLookback)
	}
	if !fromTime.Before(toTime) {
		return Range{}, fmt.Errorf("time range start must be before end")
	}

	return Range{From: fromTime, To: toTime}, nil
}

func ParseRangeWithDefault(last, from, to string, defaultLookback time.Duration, now func() time.Time) (Range, error) {
	if strings.TrimSpace(last) == "" && strings.TrimSpace(from) == "" && strings.TrimSpace(to) == "" {
		last = defaultLookback.String()
	}
	return ParseRange(last, from, to, now)
}

func parseMoment(raw string, current time.Time) (time.Time, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return time.Time{}, nil
	}
	if value == "now" {
		return current, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err == nil {
		return parsed.UTC(), nil
	}
	return time.Time{}, err
}
