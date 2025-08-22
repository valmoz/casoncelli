package casoncelli

import (
	"fmt"
	"time"
)

type AlwaysPeriod struct {
	PeriodLabel
}

func (a AlwaysPeriod) Contains(time.Time) bool {
	return true
}

func (a AlwaysPeriod) ContainsNow() bool {
	return true
}

func (a AlwaysPeriod) CurrentStart() (*time.Time, error) {
	return nil, fmt.Errorf("always period has no start")
}

func (a AlwaysPeriod) CurrentEnd() (*time.Time, error) {
	return nil, fmt.Errorf("always period has no end")
}

func (a AlwaysPeriod) NextStart() (*time.Time, error) {
	return nil, fmt.Errorf("always period has no next start")
}

func (a AlwaysPeriod) NextEnd() (*time.Time, error) {
	return nil, fmt.Errorf("always period has no next end")
}

func (a AlwaysPeriod) PreviousStart() (*time.Time, error) {
	return nil, fmt.Errorf("always period has no previous start")
}

func (a AlwaysPeriod) PreviousEnd() (*time.Time, error) {
	return nil, fmt.Errorf("always period has no previous end")
}
