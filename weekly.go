package casoncelli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type WeeklyPeriod struct {
	PeriodLabel
	From DayTimeEdge `json:"from"`
	To   DayTimeEdge `json:"to"`
}

// Contains reports whether the time instant t is included in the period.
func (p WeeklyPeriod) Contains(t time.Time) bool {
	switch {
	case p.From.Day < p.To.Day || (p.From.Day == p.To.Day && p.From.Hour < p.To.Hour):
		return p.From.BeforeOrEqual(t) && p.To.AfterOrEqual(t)
	case p.From.Day > p.To.Day || (p.From.Day == p.To.Day && p.From.Hour > p.To.Hour):
		return p.To.AfterOrEqual(t) || p.From.BeforeOrEqual(t)
	default:
		// this is the rare case where from and to are exactly the same
		return p.From.Equal(t)
	}
}

// ContainsNow reports whether the period is active.
func (p WeeklyPeriod) ContainsNow() bool {
	return p.Contains(time.Now())
}

// CurrentStart returns the start time of the current occurrance of the period, if active.
func (p WeeklyPeriod) CurrentStart() (*time.Time, error) {
	now := time.Now()
	day := now.Weekday()
	daysToRemove := 0
	if p.Contains(now) {
		if p.From.Day < day || (p.From.Day == day && p.From.Hour <= now.Format("15:04")) {
			daysToRemove = (int)(day - p.From.Day)
		} else {
			daysToRemove = (int)(7 - (p.From.Day - day))
		}
		startDay := now.AddDate(0, 0, -daysToRemove)
		startTime, err := p.From.GetEdgeTimestamp(startDay)
		return &startTime, err
	}
	return nil, fmt.Errorf("period is not active")
}

// CurrentEnd returns the end time of the current occurrance of the period, if active.
func (p WeeklyPeriod) CurrentEnd() (*time.Time, error) {
	now := time.Now()
	day := now.Weekday()
	daysToAdd := 0
	if p.Contains(now) {
		if p.To.Day > day || (p.To.Day == day && p.To.Hour >= now.Format("15:04")) {
			daysToAdd = (int)(p.To.Day - day)
		} else {
			daysToAdd = (int)(7 - (day - p.To.Day))
		}
		endDay := now.AddDate(0, 0, daysToAdd)
		endTime, err := p.To.GetEdgeTimestamp(endDay)
		return &endTime, err
	}
	return nil, fmt.Errorf("period is not active")
}

// NextStart returns the start time of the next occurrance of the period. If the period is active, it returns the start time of the next week.
func (p WeeklyPeriod) NextStart() (*time.Time, error) {
	now := time.Now()
	if p.Contains(now) {
		cs, _ := p.CurrentStart()
		ns := cs.AddDate(0, 0, 7)
		return &ns, nil
	}
	day := now.Weekday()
	daysToAdd := 0
	if p.From.Day > day || (p.From.Day == day && p.From.Hour >= now.Format("15:04")) {
		daysToAdd = (int)(p.From.Day - day)
	} else {
		daysToAdd = (int)(7 - (day - p.From.Day))
	}
	startDay := now.AddDate(0, 0, daysToAdd)
	startTime, err := p.From.GetEdgeTimestamp(startDay)
	return &startTime, err
}

// NextEnd returns the end time of the next occurrance of the period. If the period is active, it returns the end time of the next week.
func (p WeeklyPeriod) NextEnd() (*time.Time, error) {
	now := time.Now()
	if p.Contains(now) {
		ce, _ := p.CurrentEnd()
		ne := ce.AddDate(0, 0, 7)
		return &ne, nil
	}
	day := now.Weekday()
	daysToAdd := 0
	if p.To.Day > day || (p.To.Day == day && p.To.Hour >= now.Format("15:04")) {
		daysToAdd = (int)(p.To.Day - day)
	} else {
		daysToAdd = (int)(7 - (day - p.To.Day))
	}
	endDay := now.AddDate(0, 0, daysToAdd)
	endTime, err := p.To.GetEdgeTimestamp(endDay)
	return &endTime, err
}

// PreviousStart returns the start time of the previous occurrance of the period. If the period is active, it returns the start time of the previous week.
func (p WeeklyPeriod) PreviousStart() (*time.Time, error) {
	now := time.Now()
	if p.Contains(now) {
		cs, _ := p.CurrentStart()
		ps := cs.AddDate(0, 0, -7)
		return &ps, nil
	}
	day := now.Weekday()
	daysToRemove := 0
	if p.From.Day < day || (p.From.Day == day && p.From.Hour <= now.Format("15:04")) {
		daysToRemove = (int)(day - p.From.Day)
	} else {
		daysToRemove = (int)(7 - (p.From.Day - day))
	}
	startDay := now.AddDate(0, 0, -daysToRemove)
	startTime, err := p.From.GetEdgeTimestamp(startDay)
	return &startTime, err
}

// PreviousEnd returns the end time of the previous occurrance of the period. If the period is active, it returns the end time of the previous week.
func (p WeeklyPeriod) PreviousEnd() (*time.Time, error) {
	now := time.Now()
	if p.Contains(now) {
		ce, _ := p.CurrentEnd()
		pe := ce.AddDate(0, 0, -7)
		return &pe, nil
	}
	day := now.Weekday()
	daysToRemove := 0
	if p.To.Day < day || (p.To.Day == day && p.To.Hour <= now.Format("15:04")) {
		daysToRemove = (int)(day - p.To.Day)
	} else {
		daysToRemove = (int)(7 - (p.To.Day - day))
	}
	endDay := now.AddDate(0, 0, -daysToRemove)
	endTime, err := p.To.GetEdgeTimestamp(endDay)
	return &endTime, err
}

type DayTimeEdge struct {
	Day  time.Weekday `json:"day"`
	Hour string       `json:"hour"`
}

func (d *DayTimeEdge) UnmarshalJSON(data []byte) error {
	aux := struct {
		Day  string `json:"day"`
		Hour string `json:"hour"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch strings.ToLower(aux.Day) {
	case "sunday":
		d.Day = time.Sunday
	case "monday":
		d.Day = time.Monday
	case "tuesday":
		d.Day = time.Tuesday
	case "wednesday":
		d.Day = time.Wednesday
	case "thursday":
		d.Day = time.Thursday
	case "friday":
		d.Day = time.Friday
	case "saturday":
		d.Day = time.Saturday
	default:
		return fmt.Errorf("invalid weekday: %s", aux.Day)
	}
	d.Hour = aux.Hour
	return nil
}

// Before reports whether the edge is before the time instant t.
func (e DayTimeEdge) Before(t time.Time) bool {
	if e.Day < t.Weekday() {
		return true
	}
	if e.Day > t.Weekday() {
		return false
	}
	// check the hour
	edgeTimestamp, _ := e.GetEdgeTimestamp(t)
	return edgeTimestamp.Before(t)
}

// After reports whether the edge is after the time instant t.
func (e DayTimeEdge) After(t time.Time) bool {
	if e.Day > t.Weekday() {
		return true
	}
	if e.Day < t.Weekday() {
		return false
	}
	// check the hour
	edgeTimestamp, _ := e.GetEdgeTimestamp(t)
	return edgeTimestamp.After(t)
}

// Equal reports whether the edge is at the time instant t.
func (e DayTimeEdge) Equal(t time.Time) bool {
	if e.Day == t.Weekday() {
		// check the hour
		edgeTimestamp, _ := e.GetEdgeTimestamp(t)
		return edgeTimestamp.Equal(t)
	}
	return false
}

// BeforeOrEqual reports whether the edge is before or equal the time instant t.
func (e DayTimeEdge) BeforeOrEqual(t time.Time) bool {
	return e.Before(t) || e.Equal(t)
}

// AfterOrEqual reports whether the edge is after or equal the time instant t.
func (e DayTimeEdge) AfterOrEqual(t time.Time) bool {
	return e.After(t) || e.Equal(t)
}

func (e DayTimeEdge) GetEdgeTimestamp(t time.Time) (time.Time, error) {
	if t.Weekday() != e.Day {
		return time.Time{}, fmt.Errorf("day mismatch")
	}

	// Assign the hour and minute from the edge to the timestamp date
	tokens := strings.Split(e.Hour, ":")
	if len(tokens) != 2 {
		return time.Time{}, fmt.Errorf("invalid time format")
	}

	hour, err := strconv.Atoi(tokens[0])
	if err != nil {
		return time.Time{}, err
	}

	min, err := strconv.Atoi(tokens[1])
	if err != nil {
		return time.Time{}, err
	}

	edgeTimestamp := time.Date(t.Year(), t.Month(), t.Day(), hour, min, 0, 0, t.Location())
	return edgeTimestamp, nil
}
