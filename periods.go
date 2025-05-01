package casoncelli

import (
	"time"
)

type Period interface {
	Contains(time.Time) bool
	ContainsNow() bool
	CurrentStart() (*time.Time, error)
	CurrentEnd() (*time.Time, error)
	NextStart() (*time.Time, error)
	NextEnd() (*time.Time, error)
	PreviousStart() (*time.Time, error)
	PreviousEnd() (*time.Time, error)
}

type PeriodLabel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Edge interface {
	Before(time.Time) bool
	After(time.Time) bool
	Equal(time.Time) bool
	BeforeOrEqual(time.Time) bool
	AfterOrEqual(time.Time) bool
}
