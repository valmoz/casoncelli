package casoncelli

import (
	"encoding/json"
	"fmt"
	"time"
)

type OncePeriod struct {
	PeriodLabel
	From TimestampEdge `json:"from"`
	To   TimestampEdge `json:"to"`
}

func (p OncePeriod) Contains(t time.Time) bool {
	return p.From.BeforeOrEqual(t) && p.To.AfterOrEqual(t)
}

func (p OncePeriod) ContainsNow() bool {
	return p.Contains(time.Now())
}

func (p OncePeriod) CurrentStart() (*time.Time, error) {
	if p.ContainsNow() {
		return &p.From.Timestamp, nil
	}
	return nil, fmt.Errorf("period is not active")
}

func (p OncePeriod) CurrentEnd() (*time.Time, error) {
	if p.ContainsNow() {
		return &p.To.Timestamp, nil
	}
	return nil, fmt.Errorf("period is not active")
}

func (p OncePeriod) NextStart() (*time.Time, error) {
	now := time.Now()
	if p.From.After(now) {
		return &p.From.Timestamp, nil
	}
	return nil, fmt.Errorf("no next occurrence for once period")
}

func (p OncePeriod) NextEnd() (*time.Time, error) {
	now := time.Now()
	if p.From.After(now) {
		return &p.To.Timestamp, nil
	}
	return nil, fmt.Errorf("no next occurrence for once period")
}

func (p OncePeriod) PreviousStart() (*time.Time, error) {
	now := time.Now()
	if p.To.Before(now) {
		return &p.From.Timestamp, nil
	}
	return nil, fmt.Errorf("no previous occurrence for once period")
}

func (p OncePeriod) PreviousEnd() (*time.Time, error) {
	now := time.Now()
	if p.To.Before(now) {
		return &p.To.Timestamp, nil
	}
	return nil, fmt.Errorf("no previous occurrence for once period")
}

type TimestampEdge struct {
	Timestamp time.Time `json:"timestamp"`
}

func (t *TimestampEdge) UnmarshalJSON(data []byte) error {
	aux := struct {
		Timestamp string `json:"timestamp"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	ts, err := time.ParseInLocation("2006-01-02 15:04:05", aux.Timestamp, time.Local)

	if err != nil {
		return err
	}
	t.Timestamp = ts
	return nil
}

func (e TimestampEdge) Before(t time.Time) bool {
	return e.Timestamp.Before(t)
}

func (e TimestampEdge) After(t time.Time) bool {
	return e.Timestamp.After(t)
}

func (e TimestampEdge) Equal(t time.Time) bool {
	return e.Timestamp.Equal(t)
}

func (e TimestampEdge) BeforeOrEqual(t time.Time) bool {
	return e.Timestamp.Before(t) || e.Timestamp.Equal(t)
}

func (e TimestampEdge) AfterOrEqual(t time.Time) bool {
	return e.Timestamp.After(t) || e.Timestamp.Equal(t)
}
