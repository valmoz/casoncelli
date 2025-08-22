package casoncelli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAlwaysPeriodContains(t *testing.T) {
	period := AlwaysPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "always active",
			Description: "period that is always active",
		},
	}

	layout := "2006-01-02 15:04:05"

	// Test various times - should always return true
	testTime1, _ := time.Parse(layout, "2025-01-01 00:00:00")
	assert.True(t, period.Contains(testTime1), "Expected AlwaysPeriod to contain timestamp")
}

func TestAlwaysPeriodContainsNow(t *testing.T) {
	period := AlwaysPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "always active now",
			Description: "period that always contains current time",
		},
	}

	// Should always return true regardless of when test is run
	assert.True(t, period.ContainsNow(), "Expected AlwaysPeriod to contain current time")
}

func TestAlwaysPeriodCurrentStart(t *testing.T) {
	period := AlwaysPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "no current start",
			Description: "period with no defined current start",
		},
	}

	start, err := period.CurrentStart()
	assert.Nil(t, start, "Expected CurrentStart to return nil for AlwaysPeriod")
	assert.Error(t, err, "Expected CurrentStart to return error for AlwaysPeriod")
}

func TestAlwaysPeriodCurrentEnd(t *testing.T) {
	period := AlwaysPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "no current end",
			Description: "period with no defined current end",
		},
	}

	end, err := period.CurrentEnd()
	assert.Nil(t, end, "Expected CurrentEnd to return nil for AlwaysPeriod")
	assert.Error(t, err, "Expected CurrentEnd to return error for AlwaysPeriod")
}

func TestAlwaysPeriodNextStart(t *testing.T) {
	period := AlwaysPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "no next start",
			Description: "period with no next start",
		},
	}

	start, err := period.NextStart()
	assert.Nil(t, start, "Expected NextStart to return nil for AlwaysPeriod")
	assert.Error(t, err, "Expected NextStart to return error for AlwaysPeriod")
}

func TestAlwaysPeriodNextEnd(t *testing.T) {
	period := AlwaysPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "no next end",
			Description: "period with no next end",
		},
	}

	end, err := period.NextEnd()
	assert.Nil(t, end, "Expected NextEnd to return nil for AlwaysPeriod")
	assert.Error(t, err, "Expected NextEnd to return error for AlwaysPeriod")
}

func TestAlwaysPeriodPreviousStart(t *testing.T) {
	period := AlwaysPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "no previous start",
			Description: "period with no previous start",
		},
	}

	start, err := period.PreviousStart()
	assert.Nil(t, start, "Expected PreviousStart to return nil for AlwaysPeriod")
	assert.Error(t, err, "Expected PreviousStart to return error for AlwaysPeriod")
}

func TestAlwaysPeriodPreviousEnd(t *testing.T) {
	period := AlwaysPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "no previous end",
			Description: "period with no previous end",
		},
	}

	end, err := period.PreviousEnd()
	assert.Nil(t, end, "Expected PreviousEnd to return nil for AlwaysPeriod")
	assert.Error(t, err, "Expected PreviousEnd to return error for AlwaysPeriod")
}
