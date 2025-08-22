package casoncelli

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DailyPeriod struct {
	PeriodLabel
	From TimeEdge `json:"from"`
	To   TimeEdge `json:"to"`
}

// Contains reports whether the time instant t is included in the period.
func (p DailyPeriod) Contains(t time.Time) bool {
	switch {
	case p.From.Hour < p.To.Hour:
		return p.From.BeforeOrEqual(t) && p.To.AfterOrEqual(t)
	case p.From.Hour > p.To.Hour:
		return p.To.AfterOrEqual(t) || p.From.BeforeOrEqual(t)
	default:
		return p.From.Equal(t)
	}
}

// ContainsNow reports whether the period is active.
func (p DailyPeriod) ContainsNow() bool {
	return p.Contains(time.Now())
}

// CurrentStart returns the start time of the current occurrence of the period, if active.
func (p DailyPeriod) CurrentStart() (*time.Time, error) {
	now := time.Now()
	if p.Contains(now) {
		startTime, err := p.From.GetEdgeTimestamp(now)
		if p.From.Hour > p.To.Hour && now.Format("15:04") < p.From.Hour {
			startTime = startTime.AddDate(0, 0, -1)
		}
		return &startTime, err
	}
	return nil, fmt.Errorf("period is not active")
}

// CurrentEnd returns the end time of the current occurrence of the period, if active.
func (p DailyPeriod) CurrentEnd() (*time.Time, error) {
	now := time.Now()
	if p.Contains(now) {
		endTime, err := p.To.GetEdgeTimestamp(now)
		if p.From.Hour > p.To.Hour && now.Format("15:04") >= p.From.Hour {
			endTime = endTime.AddDate(0, 0, 1)
		}
		return &endTime, err
	}
	return nil, fmt.Errorf("period is not active")
}

// NextStart returns the start time of the next occurrence of the period.
func (p DailyPeriod) NextStart() (*time.Time, error) {
	now := time.Now()
	if p.Contains(now) {
		cs, _ := p.CurrentStart()
		ns := cs.AddDate(0, 0, 1)
		return &ns, nil
	}

	daysToAdd := 0
	if p.From.Hour >= now.Format("15:04") {
		daysToAdd = 0
	} else {
		daysToAdd = 1
	}

	startDay := now.AddDate(0, 0, daysToAdd)
	startTime, err := p.From.GetEdgeTimestamp(startDay)
	return &startTime, err
}

// NextEnd returns the end time of the next occurrence of the period.
func (p DailyPeriod) NextEnd() (*time.Time, error) {
	now := time.Now()
	if p.Contains(now) {
		ce, _ := p.CurrentEnd()
		ne := ce.AddDate(0, 0, 1)
		return &ne, nil
	}

	daysToAdd := 0
	if p.To.Hour >= now.Format("15:04") {
		daysToAdd = 0
	} else {
		daysToAdd = 1
	}

	endDay := now.AddDate(0, 0, daysToAdd)
	endTime, err := p.To.GetEdgeTimestamp(endDay)
	return &endTime, err
}

// PreviousStart returns the start time of the previous occurrence of the period.
func (p DailyPeriod) PreviousStart() (*time.Time, error) {
	now := time.Now()
	if p.Contains(now) {
		cs, _ := p.CurrentStart()
		ps := cs.AddDate(0, 0, -1)
		return &ps, nil
	}

	daysToRemove := 0
	if p.From.Hour <= now.Format("15:04") {
		daysToRemove = 0
	} else {
		daysToRemove = 1
	}

	startDay := now.AddDate(0, 0, -daysToRemove)
	startTime, err := p.From.GetEdgeTimestamp(startDay)
	return &startTime, err
}

// PreviousEnd returns the end time of the previous occurrence of the period.
func (p DailyPeriod) PreviousEnd() (*time.Time, error) {
	now := time.Now()
	if p.Contains(now) {
		ce, _ := p.CurrentEnd()
		pe := ce.AddDate(0, 0, -1)
		return &pe, nil
	}

	daysToRemove := 0
	if p.To.Hour <= now.Format("15:04") {
		daysToRemove = 0
	} else {
		daysToRemove = 1
	}

	endDay := now.AddDate(0, 0, -daysToRemove)
	endTime, err := p.To.GetEdgeTimestamp(endDay)
	return &endTime, err
}

type TimeEdge struct {
	Hour string `json:"hour"`
}

// Before reports whether the edge is before the time instant t.
func (e TimeEdge) Before(t time.Time) bool {
	edgeTimestamp, err := e.GetEdgeTimestamp(t)
	if err != nil {
		return false
	}
	return edgeTimestamp.Before(t)
}

// After reports whether the edge is after the time instant t.
func (e TimeEdge) After(t time.Time) bool {
	edgeTimestamp, err := e.GetEdgeTimestamp(t)
	if err != nil {
		return false
	}
	return edgeTimestamp.After(t)
}

// Equal reports whether the edge is at the time instant t.
func (e TimeEdge) Equal(t time.Time) bool {
	edgeTimestamp, err := e.GetEdgeTimestamp(t)
	if err != nil {
		return false
	}
	return edgeTimestamp.Equal(t)
}

// BeforeOrEqual reports whether the edge is before or equal the time instant t.
func (e TimeEdge) BeforeOrEqual(t time.Time) bool {
	return e.Before(t) || e.Equal(t)
}

// AfterOrEqual reports whether the edge is after or equal the time instant t.
func (e TimeEdge) AfterOrEqual(t time.Time) bool {
	return e.After(t) || e.Equal(t)
}

func (e TimeEdge) GetEdgeTimestamp(baseTime time.Time) (time.Time, error) {
	tokens := strings.Split(e.Hour, ":")
	if len(tokens) != 2 {
		return time.Time{}, fmt.Errorf("invalid hour format: %s", e.Hour)
	}

	hour, err := strconv.Atoi(tokens[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour value: %s", tokens[0])
	}

	min, err := strconv.Atoi(tokens[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute value: %s", tokens[1])
	}

	return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), hour, min, 0, 0, baseTime.Location()), nil
}
