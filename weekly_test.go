package casoncelli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWeeklyPeriodInternalContains(t *testing.T) {
	layout := "2006-01-02 15:04:05"

	// internal case test
	// tuesday 05:35 < x < thursday 22:22
	period := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  time.Tuesday,
			Hour: "07:35",
		},
		To: DayTimeEdge{
			Day:  time.Thursday,
			Hour: "22:22",
		},
	}

	ts1, _ := time.Parse(layout, "2025-04-30 08:00:00")
	assert.True(t, period.Contains(ts1), "Expected period to contain the timestamp 1 - contained case")

	ts2, _ := time.Parse(layout, "2025-04-27 07:00:00")
	assert.False(t, period.Contains(ts2), "Expected period to not contain the timestamp 2 - excluded left case")

	ts3, _ := time.Parse(layout, "2025-05-02 07:00:00")
	assert.False(t, period.Contains(ts3), "Expected period to not contain the timestamp 3 - excluded right case")

	ts4, _ := time.Parse(layout, "2025-04-23 07:00:00")
	assert.True(t, period.Contains(ts4), "Expected period to contain the timestamp 4 - previous week contained case")

	ts5, _ := time.Parse(layout, "2025-05-08 21:00:00")
	assert.True(t, period.Contains(ts5), "Expected period to contain the timestamp 5 - next week contained case")

	ts6, _ := time.Parse(layout, "2025-04-21 07:35:00")
	assert.False(t, period.Contains(ts6), "Expected period to not contain the timestamp 6 - previous week excluded left case")

	ts7, _ := time.Parse(layout, "2025-05-09 07:00:00")
	assert.False(t, period.Contains(ts7), "Expected period to not contain the timestamp 7 - next week excluded right case")

	ts8, _ := time.Parse(layout, "2025-04-29 07:34:59")
	assert.False(t, period.Contains(ts8), "Expected period to not contain the timestamp 8 - 1s excluded left case")

	ts9, _ := time.Parse(layout, "2025-04-29 07:35:00")
	assert.True(t, period.Contains(ts9), "Expected period to contain the timestamp 9 - left edge case")

	ts10, _ := time.Parse(layout, "2025-04-29 07:35:01")
	assert.True(t, period.Contains(ts10), "Expected period to contain the timestamp 10 - 1s included left case")

	ts11, _ := time.Parse(layout, "2025-05-01 22:21:59")
	assert.True(t, period.Contains(ts11), "Expected period to contain the timestamp 11 - 1s included right case")

	ts12, _ := time.Parse(layout, "2025-05-01 22:22:00")
	assert.True(t, period.Contains(ts12), "Expected period to contain the timestamp 12 - right edge case")

	ts13, _ := time.Parse(layout, "2025-05-01 22:22:01")
	assert.False(t, period.Contains(ts13), "Expected period to not contain the timestamp 13 - 1s excluded right case")
}

func TestWeeklyPeriodInternalSameDayContains(t *testing.T) {
	layout := "2006-01-02 15:04:05"

	// same day internal case test
	// wednesday 09:00 <= x <= wednesday 18:00
	period := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  time.Wednesday,
			Hour: "09:00",
		},
		To: DayTimeEdge{
			Day:  time.Wednesday,
			Hour: "18:00",
		},
	}

	ts1, _ := time.Parse(layout, "2025-04-30 10:00:00")
	assert.True(t, period.Contains(ts1), "Expected period to contain the timestamp 1 - contained case")

	ts2, _ := time.Parse(layout, "2025-04-27 07:00:00")
	assert.False(t, period.Contains(ts2), "Expected period to not contain the timestamp 2 - excluded left case")

	ts3, _ := time.Parse(layout, "2025-05-02 07:00:00")
	assert.False(t, period.Contains(ts3), "Expected period to not contain the timestamp 3 - excluded right case")

	ts4, _ := time.Parse(layout, "2025-04-23 11:00:00")
	assert.True(t, period.Contains(ts4), "Expected period to contain the timestamp 4 - previous week contained case")

	ts5, _ := time.Parse(layout, "2025-05-07 16:30:00")
	assert.True(t, period.Contains(ts5), "Expected period to contain the timestamp 5 - next week contained case")

	ts6, _ := time.Parse(layout, "2025-04-21 07:35:00")
	assert.False(t, period.Contains(ts6), "Expected period to not contain the timestamp 6 - previous week excluded left case")

	ts7, _ := time.Parse(layout, "2025-05-09 07:00:00")
	assert.False(t, period.Contains(ts7), "Expected period to not contain the timestamp 7 - next week excluded right case")

	ts8, _ := time.Parse(layout, "2025-05-30 08:59:59")
	assert.False(t, period.Contains(ts8), "Expected period to not contain the timestamp 8 - 1s excluded left case")

	ts9, _ := time.Parse(layout, "2025-04-30 09:00:00")
	assert.True(t, period.Contains(ts9), "Expected period to contain the timestamp 9 - left edge case")

	ts10, _ := time.Parse(layout, "2025-04-30 09:00:01")
	assert.True(t, period.Contains(ts10), "Expected period to contain the timestamp 10 - 1s included left case")

	ts11, _ := time.Parse(layout, "2025-04-30 17:59:59")
	assert.True(t, period.Contains(ts11), "Expected period to contain the timestamp 11 - 1s included right case")

	ts12, _ := time.Parse(layout, "2025-04-30 18:00:00")
	assert.True(t, period.Contains(ts12), "Expected period to contain the timestamp 12 - right edge case")

	ts13, _ := time.Parse(layout, "2025-04-30 18:00:01")
	assert.False(t, period.Contains(ts13), "Expected period to not contain the timestamp 13 - 1s excluded right case")
}

func TestWeeklyPeriodExternalContains(t *testing.T) {
	layout := "2006-01-02 15:04:05"

	// external case test
	// x <= tuesday 05:35 || thursday 22:22 >= x
	// equals to: thursday 22:22 <= x <= tuesday 05:35
	period := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  time.Thursday,
			Hour: "22:22",
		},
		To: DayTimeEdge{
			Day:  time.Tuesday,
			Hour: "07:35",
		},
	}

	ts1, _ := time.Parse(layout, "2025-04-30 08:00:00")
	assert.False(t, period.Contains(ts1), "Expected period to not contain the timestamp 1 - excluded case")

	ts2, _ := time.Parse(layout, "2025-04-27 07:00:00")
	assert.True(t, period.Contains(ts2), "Expected period to contain the timestamp 2 - included left case")

	ts3, _ := time.Parse(layout, "2025-05-02 07:00:00")
	assert.True(t, period.Contains(ts3), "Expected period to contain the timestamp 3 - included right case")

	ts4, _ := time.Parse(layout, "2025-04-23 07:00:00")
	assert.False(t, period.Contains(ts4), "Expected period to not contain the timestamp 4 - previous week excluded case")

	ts5, _ := time.Parse(layout, "2025-05-08 21:00:00")
	assert.False(t, period.Contains(ts5), "Expected period to not contain the timestamp 5 - next week excluded case")

	ts6, _ := time.Parse(layout, "2025-04-21 07:35:00")
	assert.True(t, period.Contains(ts6), "Expected period to contain the timestamp 6 - previous week included left case")

	ts7, _ := time.Parse(layout, "2025-05-09 07:00:00")
	assert.True(t, period.Contains(ts7), "Expected period to contain the timestamp 7 - next week included right case")

	ts8, _ := time.Parse(layout, "2025-04-29 07:34:59")
	assert.True(t, period.Contains(ts8), "Expected period to contain the timestamp 8 - 1s included left case")

	ts9, _ := time.Parse(layout, "2025-04-29 07:35:00")
	assert.True(t, period.Contains(ts9), "Expected period to contain the timestamp 9 - left edge case")

	ts10, _ := time.Parse(layout, "2025-04-29 07:35:01")
	assert.False(t, period.Contains(ts10), "Expected period to not contain the timestamp 10 - 1s excluded left case")

	ts11, _ := time.Parse(layout, "2025-05-01 22:21:59")
	assert.False(t, period.Contains(ts11), "Expected period to not contain the timestamp 11 - 1s excluded right case")

	ts12, _ := time.Parse(layout, "2025-05-01 22:22:00")
	assert.True(t, period.Contains(ts12), "Expected period to contain the timestamp 12 - right edge case")

	ts13, _ := time.Parse(layout, "2025-05-01 22:22:01")
	assert.True(t, period.Contains(ts13), "Expected period to contain the timestamp 13 - 1s included right case")
}

func TestWeeklyPeriodExternalSameDayContains(t *testing.T) {
	layout := "2006-01-02 15:04:05"

	// same day internal case test
	// wednesday 09:00 <= x <= wednesday 18:00
	period := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  time.Wednesday,
			Hour: "18:00",
		},
		To: DayTimeEdge{
			Day:  time.Wednesday,
			Hour: "09:00",
		},
	}

	ts1, _ := time.Parse(layout, "2025-04-30 10:00:00")
	assert.False(t, period.Contains(ts1), "Expected period to not contain the timestamp 1 - excluded case")

	ts2, _ := time.Parse(layout, "2025-04-27 07:00:00")
	assert.True(t, period.Contains(ts2), "Expected period to contain the timestamp 2 - included left case")

	ts3, _ := time.Parse(layout, "2025-05-02 07:00:00")
	assert.True(t, period.Contains(ts3), "Expected period to contain the timestamp 3 - included right case")

	ts4, _ := time.Parse(layout, "2025-04-23 11:00:00")
	assert.False(t, period.Contains(ts4), "Expected period to not contain the timestamp 4 - previous week excluded case")

	ts5, _ := time.Parse(layout, "2025-05-07 16:30:00")
	assert.False(t, period.Contains(ts5), "Expected period to not contain the timestamp 5 - next week excluded case")

	ts6, _ := time.Parse(layout, "2025-04-21 07:35:00")
	assert.True(t, period.Contains(ts6), "Expected period to contain the timestamp 6 - previous week included left case")

	ts7, _ := time.Parse(layout, "2025-05-09 07:00:00")
	assert.True(t, period.Contains(ts7), "Expected period to contain the timestamp 7 - next week included right case")

	ts8, _ := time.Parse(layout, "2025-05-30 08:59:59")
	assert.True(t, period.Contains(ts8), "Expected period to contain the timestamp 8 - 1s included left case")

	ts9, _ := time.Parse(layout, "2025-04-30 09:00:00")
	assert.True(t, period.Contains(ts9), "Expected period to contain the timestamp 9 - left edge case")

	ts10, _ := time.Parse(layout, "2025-04-30 09:00:01")
	assert.False(t, period.Contains(ts10), "Expected period to not contain the timestamp 10 - 1s excluded left case")

	ts11, _ := time.Parse(layout, "2025-04-30 17:59:59")
	assert.False(t, period.Contains(ts11), "Expected period to not contain the timestamp 11 - 1s excluded right case")

	ts12, _ := time.Parse(layout, "2025-04-30 18:00:00")
	assert.True(t, period.Contains(ts12), "Expected period to contain the timestamp 12 - right edge case")

	ts13, _ := time.Parse(layout, "2025-04-30 18:00:01")
	assert.True(t, period.Contains(ts13), "Expected period to contain the timestamp 13 - 1s included right case")
}

func TestWeeklyPeriodSameEdgeContains(t *testing.T) {
	layout := "2006-01-02 15:04:05"

	// same edge case test, valid only exactly in that moment.....
	// wednesday 18:00 <= x <= wednesday 18:00
	edge := DayTimeEdge{
		Day:  time.Wednesday,
		Hour: "18:00",
	}

	period := WeeklyPeriod{
		From: edge,
		To:   edge,
	}

	ts1, _ := time.Parse(layout, "2025-04-30 17:59:59")
	assert.False(t, period.Contains(ts1), "Expected period to not contain the timestamp 1 - 1s excluded left case")

	ts2, _ := time.Parse(layout, "2025-04-30 18:00:00")
	assert.True(t, period.Contains(ts2), "Expected period to contain the timestamp 2 - edge case")

	ts3, _ := time.Parse(layout, "2025-04-30 18:00:01")
	assert.False(t, period.Contains(ts3), "Expected period to not contain the timestamp 3 - 1s excluded right case")

	ts4, _ := time.Parse(layout, "2025-05-07 18:00:00")
	assert.True(t, period.Contains(ts4), "Expected period to contain the timestamp 4 - next week edge case")

	ts5, _ := time.Parse(layout, "2025-04-30 18:00:00")
	assert.True(t, period.Contains(ts5), "Expected period to contain the timestamp 5 - previous week edge case")
}

func TestWeeklyPeriodContainsNow(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)
	futureDay := future.Weekday()
	futureHour := future.Format("15:04")
	future2Day := future2.Weekday()
	future2Hour := future2.Format("15:04")
	pastDay := past.Weekday()
	pastHour := past.Format("15:04")
	past2Day := past2.Weekday()
	past2Hour := past2.Format("15:04")

	period := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
		To: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
	}

	assert.True(t, period.ContainsNow(), "Expected period to contain now - internal case")

	period2 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
		To: DayTimeEdge{
			Day:  future2Day,
			Hour: future2Hour,
		},
	}

	assert.False(t, period2.ContainsNow(), "Expected period2 to not contain now - past case")

	period3 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  past2Day,
			Hour: past2Hour,
		},
		To: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
	}

	assert.False(t, period3.ContainsNow(), "Expected period3 to not contain now - future case")

	period4 := WeeklyPeriod{
		From: period.To,
		To:   period.From,
	}

	assert.False(t, period4.ContainsNow(), "Expected period4 to not contain now - external case")

	period5 := WeeklyPeriod{
		From: period2.To,
		To:   period2.From,
	}

	assert.True(t, period5.ContainsNow(), "Expected period5 to contain now - external case right")

	period6 := WeeklyPeriod{
		From: period3.To,
		To:   period3.From,
	}

	assert.True(t, period6.ContainsNow(), "Expected period6 to contain now - external case left")
}

func TestWeeklyPeriodCurrentStart(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)
	futureDay := future.Weekday()
	futureHour := future.Format("15:04")
	future2Day := future2.Weekday()
	future2Hour := future2.Format("15:04")
	pastDay := past.Weekday()
	pastHour := past.Format("15:04")
	past2Day := past2.Weekday()
	past2Hour := past2.Format("15:04")

	period := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
		To: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
	}

	exp := time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location())
	cs, err := period.CurrentStart()
	assert.True(t, cs.Equal(exp), "Expected Current Start to match expected value for active period")
	assert.Nil(t, err, "Expected no error on current start for active period")

	period2 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
		To: DayTimeEdge{
			Day:  future2Day,
			Hour: future2Hour,
		},
	}

	cs, err = period2.CurrentStart()
	assert.Nil(t, cs, "Expected Current Start to be empty for future period")
	assert.NotNil(t, err, "Expected error on current start for future period")

	period3 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  past2Day,
			Hour: past2Hour,
		},
		To: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
	}

	cs, err = period3.CurrentStart()
	assert.Nil(t, cs, "Expected Current Start to be empty for past period")
	assert.NotNil(t, err, "Expected error on current start for past period")

	period4 := WeeklyPeriod{
		From: period.To,
		To:   period.From,
	}

	cs, err = period4.CurrentStart()
	assert.Nil(t, cs, "Expected Current Start to be empty for external period")
	assert.NotNil(t, err, "Expected error on current start for external period")

	period5 := WeeklyPeriod{
		From: period2.To,
		To:   period2.From,
	}

	//future2 is the start of the next period, so we need to subtract 7 days
	exp = time.Date(future2.Year(), future2.Month(), future2.Day(), future2.Hour(), future2.Minute(), 0, 0, future2.Location()).AddDate(0, 0, -7)
	cs, err = period5.CurrentStart()
	assert.True(t, cs.Equal(exp), "Expected Current Start to match expected value for external right period")
	assert.Nil(t, err, "Expected no error on current start for external right period")

	period6 := WeeklyPeriod{
		From: period3.To,
		To:   period3.From,
	}

	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location())
	cs, err = period6.CurrentStart()
	assert.True(t, cs.Equal(exp), "Expected Current Start to match expected value for external left period")
	assert.Nil(t, err, "Expected no error on current start for external left period")
}

func TestWeeklyPeriodCurrentEnd(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)
	futureDay := future.Weekday()
	futureHour := future.Format("15:04")
	future2Day := future2.Weekday()
	future2Hour := future2.Format("15:04")
	pastDay := past.Weekday()
	pastHour := past.Format("15:04")
	past2Day := past2.Weekday()
	past2Hour := past2.Format("15:04")

	period := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
		To: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
	}

	exp := time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location())
	cs, err := period.CurrentEnd()
	assert.True(t, cs.Equal(exp), "Expected Current End to match expected value for active period")
	assert.Nil(t, err, "Expected no error on current end for active period")

	period2 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
		To: DayTimeEdge{
			Day:  future2Day,
			Hour: future2Hour,
		},
	}

	cs, err = period2.CurrentEnd()
	assert.Nil(t, cs, "Expected Current End to be empty for future period")
	assert.NotNil(t, err, "Expected error on current end for future period")

	period3 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  past2Day,
			Hour: past2Hour,
		},
		To: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
	}

	cs, err = period3.CurrentEnd()
	assert.Nil(t, cs, "Expected Current End to be empty for past period")
	assert.NotNil(t, err, "Expected error on current end for past period")

	period4 := WeeklyPeriod{
		From: period.To,
		To:   period.From,
	}

	cs, err = period4.CurrentEnd()
	assert.Nil(t, cs, "Expected Current End to be empty for external period")
	assert.NotNil(t, err, "Expected error on current end for external period")

	period5 := WeeklyPeriod{
		From: period2.To,
		To:   period2.From,
	}

	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location())
	cs, err = period5.CurrentEnd()
	assert.True(t, cs.Equal(exp), "Expected Current End to match expected value for external right period")
	assert.Nil(t, err, "Expected no error on current end for external right period")

	period6 := WeeklyPeriod{
		From: period3.To,
		To:   period3.From,
	}

	//past2 is the end of the previous period, so we need to add 7 days
	exp = time.Date(past2.Year(), past2.Month(), past2.Day(), past2.Hour(), past2.Minute(), 0, 0, past2.Location()).AddDate(0, 0, 7)
	cs, err = period6.CurrentEnd()
	assert.True(t, cs.Equal(exp), "Expected Current End to match expected value for external left period")
	assert.Nil(t, err, "Expected no error on current end for external left period")
}

func TestWeeklyPeriodNextStart(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)
	futureDay := future.Weekday()
	futureHour := future.Format("15:04")
	future2Day := future2.Weekday()
	future2Hour := future2.Format("15:04")
	pastDay := past.Weekday()
	pastHour := past.Format("15:04")
	past2Day := past2.Weekday()
	past2Hour := past2.Format("15:04")

	period := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
		To: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
	}

	// being in active period, next start is the start of the next week period
	exp := time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location()).AddDate(0, 0, 7)
	cs, err := period.NextStart()
	assert.True(t, cs.Equal(exp), "Expected Next Start to match expected value for active period")
	assert.Nil(t, err, "Expected no error on next start for active period")

	period2 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
		To: DayTimeEdge{
			Day:  future2Day,
			Hour: future2Hour,
		},
	}

	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location())
	cs, err = period2.NextStart()
	assert.True(t, cs.Equal(exp), "Expected Next Start to match expected value for future period")
	assert.Nil(t, err, "Expected no error on next start for future period")

	period3 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  past2Day,
			Hour: past2Hour,
		},
		To: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
	}

	// being in past period, next start is the start of the next week period
	exp = time.Date(past2.Year(), past2.Month(), past2.Day(), past2.Hour(), past2.Minute(), 0, 0, past2.Location()).AddDate(0, 0, 7)
	cs, err = period3.NextStart()
	assert.True(t, cs.Equal(exp), "Expected Next Start to match expected value for past period")
	assert.Nil(t, err, "Expected no error on next start for past period")

	period4 := WeeklyPeriod{
		From: period.To,
		To:   period.From,
	}

	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location())
	cs, err = period4.NextStart()
	assert.True(t, cs.Equal(exp), "Expected Next Start to match expected value for external period")
	assert.Nil(t, err, "Expected no error on next start for external period")

	period5 := WeeklyPeriod{
		From: period2.To,
		To:   period2.From,
	}

	exp = time.Date(future2.Year(), future2.Month(), future2.Day(), future2.Hour(), future2.Minute(), 0, 0, future2.Location())
	cs, err = period5.NextStart()
	assert.True(t, cs.Equal(exp), "Expected Next Start to match expected value for external right period")
	assert.Nil(t, err, "Expected no error on next start for external right period")

	period6 := WeeklyPeriod{
		From: period3.To,
		To:   period3.From,
	}

	// being in past period, next start is the start of the next week period
	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location()).AddDate(0, 0, 7)
	cs, err = period6.NextStart()
	assert.True(t, cs.Equal(exp), "Expected Next Start to match expected value for external left period")
	assert.Nil(t, err, "Expected no error on next start for external left period")
}

func TestWeeklyPeriodNextEnd(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)
	futureDay := future.Weekday()
	futureHour := future.Format("15:04")
	future2Day := future2.Weekday()
	future2Hour := future2.Format("15:04")
	pastDay := past.Weekday()
	pastHour := past.Format("15:04")
	past2Day := past2.Weekday()
	past2Hour := past2.Format("15:04")

	period := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
		To: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
	}

	// being in active period, next end is the end of the next week period even if period didn't end
	exp := time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location()).AddDate(0, 0, 7)
	cs, err := period.NextEnd()
	assert.True(t, cs.Equal(exp), "Expected Next End to match expected value for active period")
	assert.Nil(t, err, "Expected no error on next end for active period")

	period2 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
		To: DayTimeEdge{
			Day:  future2Day,
			Hour: future2Hour,
		},
	}

	exp = time.Date(future2.Year(), future2.Month(), future2.Day(), future2.Hour(), future2.Minute(), 0, 0, future2.Location())
	cs, err = period2.NextEnd()
	assert.True(t, cs.Equal(exp), "Expected Next End to match expected value for future period")
	assert.Nil(t, err, "Expected no error on next end for future period")

	period3 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  past2Day,
			Hour: past2Hour,
		},
		To: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
	}

	// being in past period, next end is the end of the next week period
	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location()).AddDate(0, 0, 7)
	cs, err = period3.NextEnd()
	assert.True(t, cs.Equal(exp), "Expected Next End to match expected value for past period")
	assert.Nil(t, err, "Expected no error on next end for past period")

	period4 := WeeklyPeriod{
		From: period.To,
		To:   period.From,
	}

	// being in external period, next end is the end of the next week period
	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location()).AddDate(0, 0, 7)
	cs, err = period4.NextEnd()
	assert.True(t, cs.Equal(exp), "Expected Next End to match expected value for external period")
	assert.Nil(t, err, "Expected no error on next end for external period")

	period5 := WeeklyPeriod{
		From: period2.To,
		To:   period2.From,
	}

	// being active period, next end is the end of the next week period
	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location()).AddDate(0, 0, 7)
	cs, err = period5.NextEnd()
	assert.True(t, cs.Equal(exp), "Expected Next End to match expected value for external right period")
	assert.Nil(t, err, "Expected no error on next end for external right period")

	period6 := WeeklyPeriod{
		From: period3.To,
		To:   period3.From,
	}

	// being in active period, but end of this period is next week, we need to add 14 days to past2
	exp = time.Date(past2.Year(), past2.Month(), past2.Day(), past2.Hour(), past2.Minute(), 0, 0, past2.Location()).AddDate(0, 0, 14)
	cs, err = period6.NextEnd()
	assert.True(t, cs.Equal(exp), "Expected Next End to match expected value for external left period")
	assert.Nil(t, err, "Expected no error on next end for external left period")
}

func TestWeeklyPeriodPreviousStart(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)
	futureDay := future.Weekday()
	futureHour := future.Format("15:04")
	future2Day := future2.Weekday()
	future2Hour := future2.Format("15:04")
	pastDay := past.Weekday()
	pastHour := past.Format("15:04")
	past2Day := past2.Weekday()
	past2Hour := past2.Format("15:04")

	period := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
		To: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
	}

	// being in active period, previous start is the start of the previous week period
	exp := time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location()).AddDate(0, 0, -7)
	cs, err := period.PreviousStart()
	assert.True(t, cs.Equal(exp), "Expected Previous Start to match expected value for active period")
	assert.Nil(t, err, "Expected no error on previous start for active period")

	period2 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
		To: DayTimeEdge{
			Day:  future2Day,
			Hour: future2Hour,
		},
	}

	// being in future period, previous start is the start of the previous week period
	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location()).AddDate(0, 0, -7)
	cs, err = period2.PreviousStart()
	assert.True(t, cs.Equal(exp), "Expected Previous Start to match expected value for future period")
	assert.Nil(t, err, "Expected no error on previous start for future period")

	period3 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  past2Day,
			Hour: past2Hour,
		},
		To: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
	}

	exp = time.Date(past2.Year(), past2.Month(), past2.Day(), past2.Hour(), past2.Minute(), 0, 0, past2.Location())
	cs, err = period3.PreviousStart()
	assert.True(t, cs.Equal(exp), "Expected Previous Start to match expected value for past period")
	assert.Nil(t, err, "Expected no error on previous start for past period")

	period4 := WeeklyPeriod{
		From: period.To,
		To:   period.From,
	}

	// being in external period, previous start is the start of the previous week period
	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location()).AddDate(0, 0, -7)
	cs, err = period4.PreviousStart()
	assert.True(t, cs.Equal(exp), "Expected Previous Start to match expected value for external period")
	assert.Nil(t, err, "Expected no error on previous start for external period")

	period5 := WeeklyPeriod{
		From: period2.To,
		To:   period2.From,
	}

	// being in active period, but start of this period is previous week, we need to subtract 14 days to future2
	exp = time.Date(future2.Year(), future2.Month(), future2.Day(), future2.Hour(), future2.Minute(), 0, 0, future2.Location()).AddDate(0, 0, -14)
	cs, err = period5.PreviousStart()
	assert.True(t, cs.Equal(exp), "Expected Previous Start to match expected value for external right period")
	assert.Nil(t, err, "Expected no error on previous start for external right period")

	period6 := WeeklyPeriod{
		From: period3.To,
		To:   period3.From,
	}

	// being in past period, previous start is the start of the previous week period
	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location()).AddDate(0, 0, -7)
	cs, err = period6.PreviousStart()
	assert.True(t, cs.Equal(exp), "Expected Previous Start to match expected value for external left period")
	assert.Nil(t, err, "Expected no error on previous start for external left period")
}

func TestWeeklyPeriodPreviousEnd(t *testing.T) {
	now := time.Now()

	future := now.Add(2 * time.Minute)
	future2 := now.Add(4 * time.Minute)
	past := now.Add(-2 * time.Minute)
	past2 := now.Add(-4 * time.Minute)
	futureDay := future.Weekday()
	futureHour := future.Format("15:04")
	future2Day := future2.Weekday()
	future2Hour := future2.Format("15:04")
	pastDay := past.Weekday()
	pastHour := past.Format("15:04")
	past2Day := past2.Weekday()
	past2Hour := past2.Format("15:04")

	period := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
		To: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
	}

	// being in active period, previous end is the end of the previous week period
	exp := time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location()).AddDate(0, 0, -7)
	cs, err := period.PreviousEnd()
	assert.True(t, cs.Equal(exp), "Expected Previous End to match expected value for active period")
	assert.Nil(t, err, "Expected no error on previous end for active period")

	period2 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  futureDay,
			Hour: futureHour,
		},
		To: DayTimeEdge{
			Day:  future2Day,
			Hour: future2Hour,
		},
	}

	// being in future period, previous end is the end of the previous week period
	exp = time.Date(future2.Year(), future2.Month(), future2.Day(), future2.Hour(), future2.Minute(), 0, 0, future2.Location()).AddDate(0, 0, -7)
	cs, err = period2.PreviousEnd()
	assert.True(t, cs.Equal(exp), "Expected Previous End to match expected value for future period")
	assert.Nil(t, err, "Expected no error on previous end for future period")

	period3 := WeeklyPeriod{
		From: DayTimeEdge{
			Day:  past2Day,
			Hour: past2Hour,
		},
		To: DayTimeEdge{
			Day:  pastDay,
			Hour: pastHour,
		},
	}

	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location())
	cs, err = period3.PreviousEnd()
	assert.True(t, cs.Equal(exp), "Expected Previous End to match expected value for past period")
	assert.Nil(t, err, "Expected no error on previous end for past period")

	period4 := WeeklyPeriod{
		From: period.To,
		To:   period.From,
	}

	exp = time.Date(past.Year(), past.Month(), past.Day(), past.Hour(), past.Minute(), 0, 0, past.Location())
	cs, err = period4.PreviousEnd()
	assert.True(t, cs.Equal(exp), "Expected Previous End to match expected value for external period")
	assert.Nil(t, err, "Expected no error on previous end for external period")

	period5 := WeeklyPeriod{
		From: period2.To,
		To:   period2.From,
	}

	// being active period, previous end is the end of the previous week period
	exp = time.Date(future.Year(), future.Month(), future.Day(), future.Hour(), future.Minute(), 0, 0, future.Location()).AddDate(0, 0, -7)
	cs, err = period5.PreviousEnd()
	assert.True(t, cs.Equal(exp), "Expected Previous End to match expected value for external right period")
	assert.Nil(t, err, "Expected no error on previous end for external right period")

	period6 := WeeklyPeriod{
		From: period3.To,
		To:   period3.From,
	}

	exp = time.Date(past2.Year(), past2.Month(), past2.Day(), past2.Hour(), past2.Minute(), 0, 0, past2.Location())
	cs, err = period6.PreviousEnd()
	assert.True(t, cs.Equal(exp), "Expected Previous End to match expected value for external left period")
	assert.Nil(t, err, "Expected no error on previous end for external left period")
}

func TestDayTimeEdgeBefore(t *testing.T) {
	edge := DayTimeEdge{
		Day:  time.Wednesday,
		Hour: "12:00",
	}
	layout := "2006-01-02 15:04:05"
	ts, _ := time.Parse(layout, "2025-04-30 12:00:00")
	past, _ := time.Parse(layout, "2025-04-28 12:00:00")
	future, _ := time.Parse(layout, "2025-05-01 12:00:00")
	pastTime, _ := time.Parse(layout, "2025-04-30 11:59:00")
	futureTime, _ := time.Parse(layout, "2025-04-30 12:01:00")
	pastSecond, _ := time.Parse(layout, "2025-04-30 11:59:01")
	futureSecond, _ := time.Parse(layout, "2025-04-30 12:00:01")

	assert.False(t, edge.Before(past), "Expected edge to not be before the past timestamp")
	assert.False(t, edge.Before(pastTime), "Expected edge to not be before the past time timestamp")
	assert.False(t, edge.Before(pastSecond), "Expected edge to not be before the past second timestamp")
	assert.False(t, edge.Before(ts), "Expected edge to not be before the same timestamp")
	assert.True(t, edge.Before(future), "Expected edge to be before the future timestamp")
	assert.True(t, edge.Before(futureTime), "Expected edge to be before the future time timestamp")
	assert.True(t, edge.Before(futureSecond), "Expected edge to be before the future second timestamp")

	now := time.Now()

	future1 := now.Add(2 * time.Minute)
	past1 := now.Add(-2 * time.Minute)
	futureDay := future1.Weekday()
	futureHour := future1.Format("15:04")
	nowDay := now.Weekday()
	nowHour := now.Format("15:04")
	pastDay := past1.Weekday()
	pastHour := past1.Format("15:04")

	edge = DayTimeEdge{
		Day:  pastDay,
		Hour: pastHour,
	}
	assert.True(t, edge.Before(now), "Expected edge to be before now")

	edge = DayTimeEdge{
		Day:  futureDay,
		Hour: futureHour,
	}
	assert.False(t, edge.Before(now), "Expected edge to not be before now")

	edge = DayTimeEdge{
		Day:  nowDay,
		Hour: nowHour,
	}
	correctedNow := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	assert.False(t, edge.Before(correctedNow), "Expected edge to not be before now")
}

func TestDayTimeEdgeAfter(t *testing.T) {
	edge := DayTimeEdge{
		Day:  time.Wednesday,
		Hour: "12:00",
	}
	layout := "2006-01-02 15:04:05"
	ts, _ := time.Parse(layout, "2025-04-30 12:00:00")
	past, _ := time.Parse(layout, "2025-04-28 12:00:00")
	future, _ := time.Parse(layout, "2025-05-01 12:00:00")
	pastTime, _ := time.Parse(layout, "2025-04-30 11:59:00")
	futureTime, _ := time.Parse(layout, "2025-04-30 12:01:00")
	pastSecond, _ := time.Parse(layout, "2025-04-30 11:59:01")
	futureSecond, _ := time.Parse(layout, "2025-04-30 12:00:01")

	assert.True(t, edge.After(past), "Expected edge to be after the past timestamp")
	assert.True(t, edge.After(pastTime), "Expected edge to be after the past time timestamp")
	assert.True(t, edge.After(pastSecond), "Expected edge to be after the past second timestamp")
	assert.False(t, edge.After(ts), "Expected edge to not be after the same timestamp")
	assert.False(t, edge.After(future), "Expected edge to not be after the future timestamp")
	assert.False(t, edge.After(futureTime), "Expected edge to not be after the future time timestamp")
	assert.False(t, edge.After(futureSecond), "Expected edge to not be after the future second timestamp")

	now := time.Now()

	future1 := now.Add(2 * time.Minute)
	past1 := now.Add(-2 * time.Minute)
	futureDay := future1.Weekday()
	futureHour := future1.Format("15:04")
	nowDay := now.Weekday()
	nowHour := now.Format("15:04")
	pastDay := past1.Weekday()
	pastHour := past1.Format("15:04")

	edge = DayTimeEdge{
		Day:  pastDay,
		Hour: pastHour,
	}
	assert.False(t, edge.After(now), "Expected edge to not be after now")

	edge = DayTimeEdge{
		Day:  futureDay,
		Hour: futureHour,
	}
	assert.True(t, edge.After(now), "Expected edge to be after now")

	edge = DayTimeEdge{
		Day:  nowDay,
		Hour: nowHour,
	}
	correctedNow := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	assert.False(t, edge.After(correctedNow), "Expected edge to not be after now")
}

func TestDayTimeEdgeEqual(t *testing.T) {
	edge := DayTimeEdge{
		Day:  time.Wednesday,
		Hour: "12:00",
	}
	layout := "2006-01-02 15:04:05"
	ts, _ := time.Parse(layout, "2025-04-30 12:00:00")
	past, _ := time.Parse(layout, "2025-04-28 12:00:00")
	future, _ := time.Parse(layout, "2025-05-01 12:00:00")
	pastTime, _ := time.Parse(layout, "2025-04-30 11:59:00")
	futureTime, _ := time.Parse(layout, "2025-04-30 12:01:00")
	pastSecond, _ := time.Parse(layout, "2025-04-30 11:59:01")
	futureSecond, _ := time.Parse(layout, "2025-04-30 12:00:01")

	assert.False(t, edge.Equal(past), "Expected edge to not be equal the past timestamp")
	assert.False(t, edge.Equal(pastTime), "Expected edge to not be equal the past time timestamp")
	assert.False(t, edge.Equal(pastSecond), "Expected edge to not be equal the past second timestamp")
	assert.True(t, edge.Equal(ts), "Expected edge to be equal the same timestamp")
	assert.False(t, edge.Equal(future), "Expected edge to not be equal the future timestamp")
	assert.False(t, edge.Equal(futureTime), "Expected edge to not be equal the future time timestamp")
	assert.False(t, edge.Equal(futureSecond), "Expected edge to not be equal the future second timestamp")

	now := time.Now()

	future1 := now.Add(2 * time.Minute)
	past1 := now.Add(-2 * time.Minute)
	futureDay := future1.Weekday()
	futureHour := future1.Format("15:04")
	nowDay := now.Weekday()
	nowHour := now.Format("15:04")
	pastDay := past1.Weekday()
	pastHour := past1.Format("15:04")

	edge = DayTimeEdge{
		Day:  pastDay,
		Hour: pastHour,
	}
	assert.False(t, edge.Equal(now), "Expected edge to not be equal now")

	edge = DayTimeEdge{
		Day:  futureDay,
		Hour: futureHour,
	}
	assert.False(t, edge.Equal(now), "Expected edge to not be equal now")

	edge = DayTimeEdge{
		Day:  nowDay,
		Hour: nowHour,
	}
	correctedNow := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	assert.True(t, edge.Equal(correctedNow), "Expected edge to be equal now")
}

func TestDayTimeEdgeBeforeOrEqual(t *testing.T) {
	edge := DayTimeEdge{
		Day:  time.Wednesday,
		Hour: "12:00",
	}
	layout := "2006-01-02 15:04:05"
	ts, _ := time.Parse(layout, "2025-04-30 12:00:00")
	past, _ := time.Parse(layout, "2025-04-28 12:00:00")
	future, _ := time.Parse(layout, "2025-05-01 12:00:00")
	pastTime, _ := time.Parse(layout, "2025-04-30 11:59:00")
	futureTime, _ := time.Parse(layout, "2025-04-30 12:01:00")
	pastSecond, _ := time.Parse(layout, "2025-04-30 11:59:01")
	futureSecond, _ := time.Parse(layout, "2025-04-30 12:00:01")

	assert.False(t, edge.BeforeOrEqual(past), "Expected edge to not be before or equal the past timestamp")
	assert.False(t, edge.BeforeOrEqual(pastTime), "Expected edge to not be before or equal the past time timestamp")
	assert.False(t, edge.BeforeOrEqual(pastSecond), "Expected edge to not be before or equal the past second timestamp")
	assert.True(t, edge.BeforeOrEqual(ts), "Expected edge to be before or equal the same timestamp")
	assert.True(t, edge.BeforeOrEqual(future), "Expected edge to be before or equal the future timestamp")
	assert.True(t, edge.BeforeOrEqual(futureTime), "Expected edge to be before or equal the future time timestamp")
	assert.True(t, edge.BeforeOrEqual(futureSecond), "Expected edge to be before or equal the future second timestamp")

	now := time.Now()

	future1 := now.Add(2 * time.Minute)
	past1 := now.Add(-2 * time.Minute)
	futureDay := future1.Weekday()
	futureHour := future1.Format("15:04")
	nowDay := now.Weekday()
	nowHour := now.Format("15:04")
	pastDay := past1.Weekday()
	pastHour := past1.Format("15:04")

	edge = DayTimeEdge{
		Day:  pastDay,
		Hour: pastHour,
	}
	assert.True(t, edge.BeforeOrEqual(now), "Expected edge to be before or equal now")

	edge = DayTimeEdge{
		Day:  futureDay,
		Hour: futureHour,
	}
	assert.False(t, edge.BeforeOrEqual(now), "Expected edge to not be before or equal now")

	edge = DayTimeEdge{
		Day:  nowDay,
		Hour: nowHour,
	}
	correctedNow := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	assert.True(t, edge.BeforeOrEqual(correctedNow), "Expected edge to be before now")
}

func TestDayTimeEdgeAfterOrEqual(t *testing.T) {
	edge := DayTimeEdge{
		Day:  time.Wednesday,
		Hour: "12:00",
	}
	layout := "2006-01-02 15:04:05"
	ts, _ := time.Parse(layout, "2025-04-30 12:00:00")
	past, _ := time.Parse(layout, "2025-04-28 12:00:00")
	future, _ := time.Parse(layout, "2025-05-01 12:00:00")
	pastTime, _ := time.Parse(layout, "2025-04-30 11:59:00")
	futureTime, _ := time.Parse(layout, "2025-04-30 12:01:00")
	pastSecond, _ := time.Parse(layout, "2025-04-30 11:59:01")
	futureSecond, _ := time.Parse(layout, "2025-04-30 12:00:01")

	assert.True(t, edge.AfterOrEqual(past), "Expected edge to be after or equal the past timestamp")
	assert.True(t, edge.AfterOrEqual(pastTime), "Expected edge to be after or equal the past time timestamp")
	assert.True(t, edge.AfterOrEqual(pastSecond), "Expected edge to be after or equal the past second timestamp")
	assert.True(t, edge.AfterOrEqual(ts), "Expected edge to be after or equal the same timestamp")
	assert.False(t, edge.AfterOrEqual(future), "Expected edge to not be after or equal the future timestamp")
	assert.False(t, edge.AfterOrEqual(futureTime), "Expected edge to not be after or equal the future time timestamp")
	assert.False(t, edge.AfterOrEqual(futureSecond), "Expected edge to not be after or equal the future second timestamp")

	now := time.Now()

	future1 := now.Add(2 * time.Minute)
	past1 := now.Add(-2 * time.Minute)
	futureDay := future1.Weekday()
	futureHour := future1.Format("15:04")
	nowDay := now.Weekday()
	nowHour := now.Format("15:04")
	pastDay := past1.Weekday()
	pastHour := past1.Format("15:04")

	edge = DayTimeEdge{
		Day:  pastDay,
		Hour: pastHour,
	}
	assert.False(t, edge.AfterOrEqual(now), "Expected edge to not be after or equal now")

	edge = DayTimeEdge{
		Day:  futureDay,
		Hour: futureHour,
	}
	assert.True(t, edge.AfterOrEqual(now), "Expected edge to be after or equal now")

	edge = DayTimeEdge{
		Day:  nowDay,
		Hour: nowHour,
	}
	correctedNow := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	assert.True(t, edge.AfterOrEqual(correctedNow), "Expected edge to not be after or equal now")
}
