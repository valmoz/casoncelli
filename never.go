package casoncelli

import (
	"fmt"
	"time"
)

type NeverPeriod struct {
	PeriodLabel
}

func (n NeverPeriod) Contains(time.Time) bool {
	return false
}

func (n NeverPeriod) ContainsNow() bool {
	return false
}

func (n NeverPeriod) CurrentStart() (*time.Time, error) {
	return nil, fmt.Errorf("never period cannot start")
}

func (n NeverPeriod) CurrentEnd() (*time.Time, error) {
	return nil, fmt.Errorf("never period cannot end")
}

func (n NeverPeriod) NextStart() (*time.Time, error) {
	return nil, fmt.Errorf("never period has no next start")
}

func (n NeverPeriod) NextEnd() (*time.Time, error) {
	return nil, fmt.Errorf("never period has no next end")
}

func (n NeverPeriod) PreviousStart() (*time.Time, error) {
	return nil, fmt.Errorf("never period has no previous start")
}

func (n NeverPeriod) PreviousEnd() (*time.Time, error) {
	return nil, fmt.Errorf("never period has no previous end")
}
