package casoncelli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNeverPeriodContains(t *testing.T) {
	period := NeverPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "never active",
			Description: "period that is never active",
		},
	}

	layout := "2006-01-02 15:04:05"

	testTime1, _ := time.Parse(layout, "2025-01-01 00:00:00")
	assert.False(t, period.Contains(testTime1), "Expected NeverPeriod to not contain timestamp")
}

func TestNeverPeriodContainsNow(t *testing.T) {
	period := NeverPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "never active now",
			Description: "period that never contains current time",
		},
	}

	// Should always return false regardless of when test is run
	assert.False(t, period.ContainsNow(), "Expected NeverPeriod to not contain current time")
}

func TestNeverPeriodCurrentStart(t *testing.T) {
	period := NeverPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "no current start",
			Description: "period with no current start",
		},
	}

	start, err := period.CurrentStart()
	assert.Nil(t, start, "Expected CurrentStart to return nil for NeverPeriod")
	assert.Error(t, err, "Expected CurrentStart to return error for NeverPeriod")
}

func TestNeverPeriodCurrentEnd(t *testing.T) {
	period := NeverPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "no current end",
			Description: "period with no current end",
		},
	}

	end, err := period.CurrentEnd()
	assert.Nil(t, end, "Expected CurrentEnd to return nil for NeverPeriod")
	assert.Error(t, err, "Expected CurrentEnd to return error for NeverPeriod")
}

func TestNeverPeriodNextStart(t *testing.T) {
	period := NeverPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "no next start",
			Description: "period with no next start",
		},
	}

	start, err := period.NextStart()
	assert.Nil(t, start, "Expected NextStart to return nil for NeverPeriod")
	assert.Error(t, err, "Expected NextStart to return error for NeverPeriod")
}

func TestNeverPeriodNextEnd(t *testing.T) {
	period := NeverPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "no next end",
			Description: "period with no next end",
		},
	}

	end, err := period.NextEnd()
	assert.Nil(t, end, "Expected NextEnd to return nil for NeverPeriod")
	assert.Error(t, err, "Expected NextEnd to return error for NeverPeriod")
}

func TestNeverPeriodPreviousStart(t *testing.T) {
	period := NeverPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "no previous start",
			Description: "period with no previous start",
		},
	}

	start, err := period.PreviousStart()
	assert.Nil(t, start, "Expected PreviousStart to return nil for NeverPeriod")
	assert.Error(t, err, "Expected PreviousStart to return error for NeverPeriod")
}

func TestNeverPeriodPreviousEnd(t *testing.T) {
	period := NeverPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "no previous end",
			Description: "period with no previous end",
		},
	}

	end, err := period.PreviousEnd()
	assert.Nil(t, end, "Expected PreviousEnd to return nil for NeverPeriod")
	assert.Error(t, err, "Expected PreviousEnd to return error for NeverPeriod")
}
