package casoncelli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDailyPeriodInternalContains(t *testing.T) {
	// Internal case test: 09:00 <= x <= 17:00
	period := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "business hours",
			Description: "normal working hours",
		},
		From: TimeEdge{Hour: "09:00"},
		To:   TimeEdge{Hour: "17:00"},
	}

	layout := "2006-01-02 15:04:05"

	// Test time before the period
	beforeTime, _ := time.Parse(layout, "2025-08-22 08:30:00")
	assert.False(t, period.Contains(beforeTime), "Expected period to not contain time before start (08:30)")

	// Test time at start of period
	startTime, _ := time.Parse(layout, "2025-08-22 09:00:00")
	assert.True(t, period.Contains(startTime), "Expected period to contain time at start (09:00)")

	// Test time within period
	withinTime, _ := time.Parse(layout, "2025-08-22 12:30:00")
	assert.True(t, period.Contains(withinTime), "Expected period to contain time within period (12:30)")

	// Test time at end of period
	endTime, _ := time.Parse(layout, "2025-08-22 17:00:00")
	assert.True(t, period.Contains(endTime), "Expected period to contain time at end (17:00)")

	// Test time after the period
	afterTime, _ := time.Parse(layout, "2025-08-22 18:30:00")
	assert.False(t, period.Contains(afterTime), "Expected period to not contain time after end (18:30)")
}

func TestDailyPeriodSameHourContains(t *testing.T) {
	// Same hour case test: 12:00 <= x <= 12:00
	period := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "lunch break",
			Description: "exact lunch time",
		},
		From: TimeEdge{Hour: "12:00"},
		To:   TimeEdge{Hour: "12:00"},
	}

	layout := "2006-01-02 15:04:05"

	// Test time before the exact hour
	beforeTime, _ := time.Parse(layout, "2025-08-22 11:59:00")
	assert.False(t, period.Contains(beforeTime), "Expected period to not contain time before exact hour (11:59)")

	// Test exact time
	exactTime, _ := time.Parse(layout, "2025-08-22 12:00:00")
	assert.True(t, period.Contains(exactTime), "Expected period to contain exact time (12:00)")

	// Test time after the exact hour
	afterTime, _ := time.Parse(layout, "2025-08-22 12:01:00")
	assert.False(t, period.Contains(afterTime), "Expected period to not contain time after exact hour (12:01)")
}

func TestDailyPeriodExternalContains(t *testing.T) {
	// External case test: crosses midnight 22:00 <= x <= 06:00
	period := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "night shift",
			Description: "overnight period crossing midnight",
		},
		From: TimeEdge{Hour: "22:00"},
		To:   TimeEdge{Hour: "06:00"},
	}

	layout := "2006-01-02 15:04:05"

	// Test time in the evening part (after 22:00)
	eveningTime, _ := time.Parse(layout, "2025-08-22 23:30:00")
	assert.True(t, period.Contains(eveningTime), "Expected period to contain evening time (23:30)")

	// Test time at start of period
	startTime, _ := time.Parse(layout, "2025-08-22 22:00:00")
	assert.True(t, period.Contains(startTime), "Expected period to contain time at start (22:00)")

	// Test time in the morning part (before 06:00)
	morningTime, _ := time.Parse(layout, "2025-08-22 03:30:00")
	assert.True(t, period.Contains(morningTime), "Expected period to contain morning time (03:30)")

	// Test time at end of period
	endTime, _ := time.Parse(layout, "2025-08-22 06:00:00")
	assert.True(t, period.Contains(endTime), "Expected period to contain time at end (06:00)")

	// Test time in the middle of the day (excluded)
	middayTime, _ := time.Parse(layout, "2025-08-22 12:00:00")
	assert.False(t, period.Contains(middayTime), "Expected period to not contain midday time (12:00)")

	// Test time just before start
	beforeStartTime, _ := time.Parse(layout, "2025-08-22 21:59:00")
	assert.False(t, period.Contains(beforeStartTime), "Expected period to not contain time just before start (21:59)")

	// Test time just after end
	afterEndTime, _ := time.Parse(layout, "2025-08-22 06:01:00")
	assert.False(t, period.Contains(afterEndTime), "Expected period to not contain time just after end (06:01)")
}

func TestDailyPeriodContainsNow(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)

	futureHour := future.Format("15:04")
	future2Hour := future2.Format("15:04")
	pastHour := past.Format("15:04")
	past2Hour := past2.Format("15:04")

	// Test internal case: past <= now <= future
	period := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "active period",
			Description: "period containing now",
		},
		From: TimeEdge{Hour: pastHour},
		To:   TimeEdge{Hour: futureHour},
	}

	assert.True(t, period.ContainsNow(), "Expected period to contain now - internal case")

	// Test future case: both times are in the future
	period2 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "future period",
			Description: "period in the future",
		},
		From: TimeEdge{Hour: futureHour},
		To:   TimeEdge{Hour: future2Hour},
	}

	assert.False(t, period2.ContainsNow(), "Expected period2 to not contain now - future case")

	// Test past case: both times are in the past
	period3 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "past period",
			Description: "period in the past",
		},
		From: TimeEdge{Hour: past2Hour},
		To:   TimeEdge{Hour: pastHour},
	}

	assert.False(t, period3.ContainsNow(), "Expected period3 to not contain now - past case")

	// Test external case (crossing midnight): future <= now <= past (inverted)
	period4 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "external period",
			Description: "period crossing midnight, excluding now",
		},
		From: period.To,   // futureHour
		To:   period.From, // pastHour
	}

	assert.False(t, period4.ContainsNow(), "Expected period4 to not contain now - external case")

	// Test external case right (crossing midnight): future2 <= now <= future
	period5 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "external period right",
			Description: "period crossing midnight, including now on right side",
		},
		From: period2.To,   // future2Hour
		To:   period2.From, // futureHour
	}

	assert.True(t, period5.ContainsNow(), "Expected period5 to contain now - external case right")

	// Test external case left (crossing midnight): past <= now <= past2
	period6 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "external period left",
			Description: "period crossing midnight, including now on left side",
		},
		From: period3.To,   // pastHour
		To:   period3.From, // past2Hour
	}

	assert.True(t, period6.ContainsNow(), "Expected period6 to contain now - external case left")
}

func TestDailyPeriodContainsWithMinutes(t *testing.T) {
	// Test with specific minutes
	period := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "coffee break",
			Description: "morning coffee break",
		},
		From: TimeEdge{Hour: "10:15"},
		To:   TimeEdge{Hour: "10:45"},
	}

	layout := "2006-01-02 15:04:05"

	// Test time before the period
	beforeTime, _ := time.Parse(layout, "2025-08-22 10:14:59")
	assert.False(t, period.Contains(beforeTime), "Expected period to not contain time just before start (10:14:59)")

	// Test time at start of period
	startTime, _ := time.Parse(layout, "2025-08-22 10:15:00")
	assert.True(t, period.Contains(startTime), "Expected period to contain time at start (10:15:00)")

	// Test time within period
	withinTime, _ := time.Parse(layout, "2025-08-22 10:30:00")
	assert.True(t, period.Contains(withinTime), "Expected period to contain time within period (10:30:00)")

	// Test time at end of period
	endTime, _ := time.Parse(layout, "2025-08-22 10:45:00")
	assert.True(t, period.Contains(endTime), "Expected period to contain time at end (10:45:00)")

	// Test time after the period
	afterTime, _ := time.Parse(layout, "2025-08-22 10:45:01")
	assert.False(t, period.Contains(afterTime), "Expected period to not contain time just after end (10:45:01)")
}

func TestDailyPeriodCurrentStart(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)

	futureHour := future.Format("15:04")
	future2Hour := future2.Format("15:04")
	pastHour := past.Format("15:04")
	past2Hour := past2.Format("15:04")

	// Test internal case: past <= now <= future
	period := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "active period",
			Description: "period containing now",
		},
		From: TimeEdge{Hour: pastHour},
		To:   TimeEdge{Hour: futureHour},
	}

	exp := time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location())
	cs, err := period.CurrentStart()
	assert.True(t, cs.Equal(exp), "Expected Current Start to match expected value for active period")
	assert.Nil(t, err, "Expected no error on current start for active period")

	// Test future case: both times are in the future
	period2 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "future period",
			Description: "period in the future",
		},
		From: TimeEdge{Hour: futureHour},
		To:   TimeEdge{Hour: future2Hour},
	}

	cs, err = period2.CurrentStart()
	assert.Nil(t, cs, "Expected Current Start to be empty for future period")
	assert.NotNil(t, err, "Expected error on current start for future period")

	// Test past case: both times are in the past
	period3 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "past period",
			Description: "period in the past",
		},
		From: TimeEdge{Hour: past2Hour},
		To:   TimeEdge{Hour: pastHour},
	}

	cs, err = period3.CurrentStart()
	assert.Nil(t, cs, "Expected Current Start to be empty for past period")
	assert.NotNil(t, err, "Expected error on current start for past period")

	// Test external case (crossing midnight): future <= now <= past (inverted)
	period4 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "external period",
			Description: "period crossing midnight, excluding now",
		},
		From: period.To,   // futureHour
		To:   period.From, // pastHour
	}

	cs, err = period4.CurrentStart()
	assert.Nil(t, cs, "Expected Current Start to be empty for external period")
	assert.NotNil(t, err, "Expected error on current start for external period")

	// Test external case right (crossing midnight): future2 <= now <= future
	period5 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "external period right",
			Description: "period crossing midnight, including now on right side",
		},
		From: period2.To,   // future2Hour
		To:   period2.From, // futureHour
	}

	// future2 is the start of the next period, so we need to subtract 1 day
	exp = time.Date(future2.Year(), future2.Month(), future2.Day(), future2.Hour(), future2.Minute(), 0, 0, future2.Location()).AddDate(0, 0, -1)
	cs, err = period5.CurrentStart()
	assert.True(t, cs.Equal(exp), "Expected Current Start to match expected value for external right period")
	assert.Nil(t, err, "Expected no error on current start for external right period")

	// Test external case left (crossing midnight): past <= now <= past2
	period6 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "external period left",
			Description: "period crossing midnight, including now on left side",
		},
		From: period3.To,   // pastHour
		To:   period3.From, // past2Hour
	}

	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location())
	cs, err = period6.CurrentStart()
	assert.True(t, cs.Equal(exp), "Expected Current Start to match expected value for external left period")
	assert.Nil(t, err, "Expected no error on current start for external left period")
}

func TestDailyPeriodCurrentEnd(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)

	futureHour := future.Format("15:04")
	future2Hour := future2.Format("15:04")
	pastHour := past.Format("15:04")
	past2Hour := past2.Format("15:04")

	// Test internal case: past <= now <= future
	period := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "active period",
			Description: "period containing now",
		},
		From: TimeEdge{Hour: pastHour},
		To:   TimeEdge{Hour: futureHour},
	}

	exp := time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location())
	cs, err := period.CurrentEnd()
	assert.True(t, cs.Equal(exp), "Expected Current End to match expected value for active period")
	assert.Nil(t, err, "Expected no error on current end for active period")

	// Test future case: both times are in the future
	period2 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "future period",
			Description: "period in the future",
		},
		From: TimeEdge{Hour: futureHour},
		To:   TimeEdge{Hour: future2Hour},
	}

	cs, err = period2.CurrentEnd()
	assert.Nil(t, cs, "Expected Current End to be empty for future period")
	assert.NotNil(t, err, "Expected error on current end for future period")

	// Test past case: both times are in the past
	period3 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "past period",
			Description: "period in the past",
		},
		From: TimeEdge{Hour: past2Hour},
		To:   TimeEdge{Hour: pastHour},
	}

	cs, err = period3.CurrentEnd()
	assert.Nil(t, cs, "Expected Current End to be empty for past period")
	assert.NotNil(t, err, "Expected error on current end for past period")

	// Test external case (crossing midnight): future <= now <= past (inverted)
	period4 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "external period",
			Description: "period crossing midnight, excluding now",
		},
		From: period.To,   // futureHour
		To:   period.From, // pastHour
	}

	cs, err = period4.CurrentEnd()
	assert.Nil(t, cs, "Expected Current End to be empty for external period")
	assert.NotNil(t, err, "Expected error on current end for external period")

	// Test external case right (crossing midnight): future2 <= now <= future
	period5 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "external period right",
			Description: "period crossing midnight, including now on right side",
		},
		From: period2.To,   // future2Hour
		To:   period2.From, // futureHour
	}

	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location())
	cs, err = period5.CurrentEnd()
	assert.True(t, cs.Equal(exp), "Expected Current End to match expected value for external right period")
	assert.Nil(t, err, "Expected no error on current end for external right period")

	// Test external case left (crossing midnight): past <= now <= past2
	period6 := DailyPeriod{
		PeriodLabel: PeriodLabel{
			Name:        "external period left",
			Description: "period crossing midnight, including now on left side",
		},
		From: period3.To,   // pastHour
		To:   period3.From, // past2Hour
	}

	// past2 is the end of the previous period, so we need to add 1 day
	exp = time.Date(past2.Year(), past2.Month(), past2.Day(), past2.Hour(), past2.Minute(), 0, 0, past2.Location()).AddDate(0, 0, 1)
	cs, err = period6.CurrentEnd()
	assert.True(t, cs.Equal(exp), "Expected Current End to match expected value for external left period")
	assert.Nil(t, err, "Expected no error on current end for external left period")
}

func TestDailyPeriodNextStart(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)

	futureHour := future.Format("15:04")
	future2Hour := future2.Format("15:04")
	pastHour := past.Format("15:04")
	past2Hour := past2.Format("15:04")

	period := DailyPeriod{
		From: TimeEdge{Hour: pastHour},
		To:   TimeEdge{Hour: futureHour},
	}

	// being in active period, next start is the start of the next day period
	exp := time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location()).AddDate(0, 0, 1)
	cs, err := period.NextStart()
	assert.True(t, cs.Equal(exp), "Expected Next Start to match expected value for active period")
	assert.Nil(t, err, "Expected no error on next start for active period")

	period2 := DailyPeriod{
		From: TimeEdge{Hour: futureHour},
		To:   TimeEdge{Hour: future2Hour},
	}

	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location())
	cs, err = period2.NextStart()
	assert.True(t, cs.Equal(exp), "Expected Next Start to match expected value for future period")
	assert.Nil(t, err, "Expected no error on next start for future period")

	period3 := DailyPeriod{
		From: TimeEdge{Hour: past2Hour},
		To:   TimeEdge{Hour: pastHour},
	}

	// being in past period, next start is the start of the next day period
	exp = time.Date(past2.Year(), past2.Month(), past2.Day(), past2.Hour(), past2.Minute(), 0, 0, past2.Location()).AddDate(0, 0, 1)
	cs, err = period3.NextStart()
	assert.True(t, cs.Equal(exp), "Expected Next Start to match expected value for past period")
	assert.Nil(t, err, "Expected no error on next start for past period")

	period4 := DailyPeriod{
		From: period.To,
		To:   period.From,
	}

	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location())
	cs, err = period4.NextStart()
	assert.True(t, cs.Equal(exp), "Expected Next Start to match expected value for external period")
	assert.Nil(t, err, "Expected no error on next start for external period")

	period5 := DailyPeriod{
		From: period2.To,
		To:   period2.From,
	}

	exp = time.Date(future2.Year(), future2.Month(), future2.Day(), future2.Hour(), future2.Minute(), 0, 0, future2.Location())
	cs, err = period5.NextStart()
	assert.True(t, cs.Equal(exp), "Expected Next Start to match expected value for external right period")
	assert.Nil(t, err, "Expected no error on next start for external right period")

	period6 := DailyPeriod{
		From: period3.To,
		To:   period3.From,
	}

	// being in past period, next start is the start of the next day period
	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location()).AddDate(0, 0, 1)
	cs, err = period6.NextStart()
	assert.True(t, cs.Equal(exp), "Expected Next Start to match expected value for external left period")
	assert.Nil(t, err, "Expected no error on next start for external left period")
}

func TestDailyPeriodNextEnd(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)

	futureHour := future.Format("15:04")
	future2Hour := future2.Format("15:04")
	pastHour := past.Format("15:04")
	past2Hour := past2.Format("15:04")

	period := DailyPeriod{
		From: TimeEdge{Hour: pastHour},
		To:   TimeEdge{Hour: futureHour},
	}

	// being in active period, next end is the end of the next day period even if period didn't end
	exp := time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location()).AddDate(0, 0, 1)
	cs, err := period.NextEnd()
	assert.True(t, cs.Equal(exp), "Expected Next End to match expected value for active period")
	assert.Nil(t, err, "Expected no error on next end for active period")

	period2 := DailyPeriod{
		From: TimeEdge{Hour: futureHour},
		To:   TimeEdge{Hour: future2Hour},
	}

	exp = time.Date(future2.Year(), future2.Month(), future2.Day(), future2.Hour(), future2.Minute(), 0, 0, future2.Location())
	cs, err = period2.NextEnd()
	assert.True(t, cs.Equal(exp), "Expected Next End to match expected value for future period")
	assert.Nil(t, err, "Expected no error on next end for future period")

	period3 := DailyPeriod{
		From: TimeEdge{Hour: past2Hour},
		To:   TimeEdge{Hour: pastHour},
	}

	// being in past period, next end is the end of the next day period
	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location()).AddDate(0, 0, 1)
	cs, err = period3.NextEnd()
	assert.True(t, cs.Equal(exp), "Expected Next End to match expected value for past period")
	assert.Nil(t, err, "Expected no error on next end for past period")

	period4 := DailyPeriod{
		From: period.To,
		To:   period.From,
	}

	// being in external period, next end is the end of the next day period
	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location()).AddDate(0, 0, 1)
	cs, err = period4.NextEnd()
	assert.True(t, cs.Equal(exp), "Expected Next End to match expected value for external period")
	assert.Nil(t, err, "Expected no error on next end for external period")

	period5 := DailyPeriod{
		From: period2.To,
		To:   period2.From,
	}

	// being active period, next end is the end of the next day period
	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location()).AddDate(0, 0, 1)
	cs, err = period5.NextEnd()
	assert.True(t, cs.Equal(exp), "Expected Next End to match expected value for external right period")
	assert.Nil(t, err, "Expected no error on next end for external right period")

	period6 := DailyPeriod{
		From: period3.To,
		To:   period3.From,
	}

	// being in active period, but end of this period is next day, we need to add 2 days to past2
	exp = time.Date(past2.Year(), past2.Month(), past2.Day(), past2.Hour(), past2.Minute(), 0, 0, past2.Location()).AddDate(0, 0, 2)
	cs, err = period6.NextEnd()
	assert.True(t, cs.Equal(exp), "Expected Next End to match expected value for external left period")
	assert.Nil(t, err, "Expected no error on next end for external left period")
}

func TestDailyPeriodPreviousStart(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)

	futureHour := future.Format("15:04")
	future2Hour := future2.Format("15:04")
	pastHour := past.Format("15:04")
	past2Hour := past2.Format("15:04")

	period := DailyPeriod{
		From: TimeEdge{Hour: pastHour},
		To:   TimeEdge{Hour: futureHour},
	}

	// being in active period, previous start is the start of the previous day period
	exp := time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location()).AddDate(0, 0, -1)
	cs, err := period.PreviousStart()
	assert.True(t, cs.Equal(exp), "Expected Previous Start to match expected value for active period")
	assert.Nil(t, err, "Expected no error on previous start for active period")

	period2 := DailyPeriod{
		From: TimeEdge{Hour: futureHour},
		To:   TimeEdge{Hour: future2Hour},
	}

	// being in future period, previous start is the start of the previous day period
	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location()).AddDate(0, 0, -1)
	cs, err = period2.PreviousStart()
	assert.True(t, cs.Equal(exp), "Expected Previous Start to match expected value for future period")
	assert.Nil(t, err, "Expected no error on previous start for future period")

	period3 := DailyPeriod{
		From: TimeEdge{Hour: past2Hour},
		To:   TimeEdge{Hour: pastHour},
	}

	exp = time.Date(past2.Year(), past2.Month(), past2.Day(), past2.Hour(), past2.Minute(), 0, 0, past2.Location())
	cs, err = period3.PreviousStart()
	assert.True(t, cs.Equal(exp), "Expected Previous Start to match expected value for past period")
	assert.Nil(t, err, "Expected no error on previous start for past period")

	period4 := DailyPeriod{
		From: period.To,
		To:   period.From,
	}

	// being in external period, previous start is the start of the previous day period
	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location()).AddDate(0, 0, -1)
	cs, err = period4.PreviousStart()
	assert.True(t, cs.Equal(exp), "Expected Previous Start to match expected value for external period")
	assert.Nil(t, err, "Expected no error on previous start for external period")

	period5 := DailyPeriod{
		From: period2.To,
		To:   period2.From,
	}

	// being in active period, but start of this period is previous day, we need to subtract 2 days to future2
	exp = time.Date(future2.Year(), future2.Month(), future2.Day(), future2.Hour(), future2.Minute(), 0, 0, future2.Location()).AddDate(0, 0, -2)
	cs, err = period5.PreviousStart()
	assert.True(t, cs.Equal(exp), "Expected Previous Start to match expected value for external right period")
	assert.Nil(t, err, "Expected no error on previous start for external right period")

	period6 := DailyPeriod{
		From: period3.To,
		To:   period3.From,
	}

	// being in past period, previous start is the start of the previous day period
	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location()).AddDate(0, 0, -1)
	cs, err = period6.PreviousStart()
	assert.True(t, cs.Equal(exp), "Expected Previous Start to match expected value for external left period")
	assert.Nil(t, err, "Expected no error on previous start for external left period")
}

func TestDailyPeriodPreviousEnd(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)

	futureHour := future.Format("15:04")
	future2Hour := future2.Format("15:04")
	pastHour := past.Format("15:04")
	past2Hour := past2.Format("15:04")

	period := DailyPeriod{
		From: TimeEdge{Hour: pastHour},
		To:   TimeEdge{Hour: futureHour},
	}

	// being in active period, previous end is the end of the previous day period
	exp := time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location()).AddDate(0, 0, -1)
	cs, err := period.PreviousEnd()
	assert.True(t, cs.Equal(exp), "Expected Previous End to match expected value for active period")
	assert.Nil(t, err, "Expected no error on previous end for active period")

	period2 := DailyPeriod{
		From: TimeEdge{Hour: futureHour},
		To:   TimeEdge{Hour: future2Hour},
	}

	// being in future period, previous end is the end of the previous day period
	exp = time.Date(future2.Year(), future2.Month(), future2.Day(), future2.Hour(), future2.Minute(), 0, 0, future2.Location()).AddDate(0, 0, -1)
	cs, err = period2.PreviousEnd()
	assert.True(t, cs.Equal(exp), "Expected Previous End to match expected value for future period")
	assert.Nil(t, err, "Expected no error on previous end for future period")

	period3 := DailyPeriod{
		From: TimeEdge{Hour: past2Hour},
		To:   TimeEdge{Hour: pastHour},
	}

	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location())
	cs, err = period3.PreviousEnd()
	assert.True(t, cs.Equal(exp), "Expected Previous End to match expected value for past period")
	assert.Nil(t, err, "Expected no error on previous end for past period")

	period4 := DailyPeriod{
		From: period.To,
		To:   period.From,
	}

	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location())
	cs, err = period4.PreviousEnd()
	assert.True(t, cs.Equal(exp), "Expected Previous End to match expected value for external period")
	assert.Nil(t, err, "Expected no error on previous end for external period")

	period5 := DailyPeriod{
		From: period2.To,
		To:   period2.From,
	}

	// being active period, previous end is the end of the previous day period
	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location()).AddDate(0, 0, -1)
	cs, err = period5.PreviousEnd()
	assert.True(t, cs.Equal(exp), "Expected Previous End to match expected value for external right period")
	assert.Nil(t, err, "Expected no error on previous end for external right period")

	period6 := DailyPeriod{
		From: period3.To,
		To:   period3.From,
	}

	exp = time.Date(past2.Year(), past2.Month(), past2.Day(), past2.Hour(), past2.Minute(), 0, 0, past2.Location())
	cs, err = period6.PreviousEnd()
	assert.True(t, cs.Equal(exp), "Expected Previous End to match expected value for external left period")
	assert.Nil(t, err, "Expected no error on previous end for external left period")
}

func TestTimeEdgeBefore(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	baseTime, _ := time.Parse(layout, "2025-08-22 12:00:00")

	// Test edge before the time
	earlyEdge := TimeEdge{Hour: "10:00"}
	assert.True(t, earlyEdge.Before(baseTime), "Expected edge (10:00) to be before the base time (12:00)")

	// Test edge after the time
	lateEdge := TimeEdge{Hour: "14:00"}
	assert.False(t, lateEdge.Before(baseTime), "Expected edge (14:00) to not be before the base time (12:00)")

	// Test edge at the same time
	sameEdge := TimeEdge{Hour: "12:00"}
	assert.False(t, sameEdge.Before(baseTime), "Expected edge (12:00) to not be before the same base time (12:00)")

	// Test with minutes
	earlyMinuteEdge := TimeEdge{Hour: "11:59"}
	assert.True(t, earlyMinuteEdge.Before(baseTime), "Expected edge (11:59) to be before the base time (12:00)")

	lateMinuteEdge := TimeEdge{Hour: "12:01"}
	assert.False(t, lateMinuteEdge.Before(baseTime), "Expected edge (12:01) to not be before the base time (12:00)")
}

func TestTimeEdgeAfter(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	baseTime, _ := time.Parse(layout, "2025-08-22 12:00:00")

	// Test edge after the time
	lateEdge := TimeEdge{Hour: "14:00"}
	assert.True(t, lateEdge.After(baseTime), "Expected edge (14:00) to be after the base time (12:00)")

	// Test edge before the time
	earlyEdge := TimeEdge{Hour: "10:00"}
	assert.False(t, earlyEdge.After(baseTime), "Expected edge (10:00) to not be after the base time (12:00)")

	// Test edge at the same time
	sameEdge := TimeEdge{Hour: "12:00"}
	assert.False(t, sameEdge.After(baseTime), "Expected edge (12:00) to not be after the same base time (12:00)")

	// Test with minutes
	lateMinuteEdge := TimeEdge{Hour: "12:01"}
	assert.True(t, lateMinuteEdge.After(baseTime), "Expected edge (12:01) to be after the base time (12:00)")

	earlyMinuteEdge := TimeEdge{Hour: "11:59"}
	assert.False(t, earlyMinuteEdge.After(baseTime), "Expected edge (11:59) to not be after the base time (12:00)")
}

func TestTimeEdgeEqual(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	baseTime, _ := time.Parse(layout, "2025-08-22 12:00:00")

	// Test edge equal to the time
	sameEdge := TimeEdge{Hour: "12:00"}
	assert.True(t, sameEdge.Equal(baseTime), "Expected edge (12:00) to be equal to the base time (12:00)")

	// Test edge before the time
	earlyEdge := TimeEdge{Hour: "11:59"}
	assert.False(t, earlyEdge.Equal(baseTime), "Expected edge (11:59) to not be equal to the base time (12:00)")

	// Test edge after the time
	lateEdge := TimeEdge{Hour: "12:01"}
	assert.False(t, lateEdge.Equal(baseTime), "Expected edge (12:01) to not be equal to the base time (12:00)")

	// Test with different hours
	differentEdge := TimeEdge{Hour: "15:30"}
	assert.False(t, differentEdge.Equal(baseTime), "Expected edge (15:30) to not be equal to the base time (12:00)")
}

func TestTimeEdgeBeforeOrEqual(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	baseTime, _ := time.Parse(layout, "2025-08-22 12:00:00")

	// Test edge before the time
	earlyEdge := TimeEdge{Hour: "10:00"}
	assert.True(t, earlyEdge.BeforeOrEqual(baseTime), "Expected edge (10:00) to be before or equal to the base time (12:00)")

	// Test edge equal to the time
	sameEdge := TimeEdge{Hour: "12:00"}
	assert.True(t, sameEdge.BeforeOrEqual(baseTime), "Expected edge (12:00) to be before or equal to the base time (12:00)")

	// Test edge after the time
	lateEdge := TimeEdge{Hour: "14:00"}
	assert.False(t, lateEdge.BeforeOrEqual(baseTime), "Expected edge (14:00) to not be before or equal to the base time (12:00)")

	// Test with minutes - before
	earlyMinuteEdge := TimeEdge{Hour: "11:59"}
	assert.True(t, earlyMinuteEdge.BeforeOrEqual(baseTime), "Expected edge (11:59) to be before or equal to the base time (12:00)")

	// Test with minutes - after
	lateMinuteEdge := TimeEdge{Hour: "12:01"}
	assert.False(t, lateMinuteEdge.BeforeOrEqual(baseTime), "Expected edge (12:01) to not be before or equal to the base time (12:00)")
}

func TestTimeEdgeAfterOrEqual(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	baseTime, _ := time.Parse(layout, "2025-08-22 12:00:00")

	// Test edge after the time
	lateEdge := TimeEdge{Hour: "14:00"}
	assert.True(t, lateEdge.AfterOrEqual(baseTime), "Expected edge (14:00) to be after or equal to the base time (12:00)")

	// Test edge equal to the time
	sameEdge := TimeEdge{Hour: "12:00"}
	assert.True(t, sameEdge.AfterOrEqual(baseTime), "Expected edge (12:00) to be after or equal to the base time (12:00)")

	// Test edge before the time
	earlyEdge := TimeEdge{Hour: "10:00"}
	assert.False(t, earlyEdge.AfterOrEqual(baseTime), "Expected edge (10:00) to not be after or equal to the base time (12:00)")

	// Test with minutes - after
	lateMinuteEdge := TimeEdge{Hour: "12:01"}
	assert.True(t, lateMinuteEdge.AfterOrEqual(baseTime), "Expected edge (12:01) to be after or equal to the base time (12:00)")

	// Test with minutes - before
	earlyMinuteEdge := TimeEdge{Hour: "11:59"}
	assert.False(t, earlyMinuteEdge.AfterOrEqual(baseTime), "Expected edge (11:59) to not be after or equal to the base time (12:00)")
}
