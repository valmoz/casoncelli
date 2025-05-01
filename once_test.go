package casoncelli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOncePeriodContains(t *testing.T) {
	layout := "2006-01-02 15:04"
	ts, _ := time.Parse(layout, "2025-04-28 12:00")
	period := OncePeriod{
		From: TimestampEdge{Timestamp: ts},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
	}

	past, _ := time.Parse(layout, "2025-04-21 07:35")
	contained, _ := time.Parse(layout, "2025-04-28 20:00")
	future, _ := time.Parse(layout, "2025-04-29 12:01")

	assert.False(t, period.Contains(past), "Expected period to not contain the past timestamp")
	assert.True(t, period.Contains(ts), "Expected period to contain the from timestamp")
	assert.True(t, period.Contains(contained), "Expected period to contain the contained timestamp")
	assert.True(t, period.Contains(ts.AddDate(0, 0, 1)), "Expected period to contain the to timestamp")
	assert.False(t, period.Contains(future), "Expected period to not contain the future timestamp")
}

func TestOncePeriodContainsNow(t *testing.T) {
	ts := time.Now()
	period := OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
	}

	assert.True(t, period.ContainsNow(), "Expected period to contain the current timestamp")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 2)},
	}

	assert.False(t, period.ContainsNow(), "Expected future period to not contain the current timestamp")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -2)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
	}

	assert.False(t, period.ContainsNow(), "Expected past period to not contain the current timestamp")
}

func TestOncePeriodCurrentStart(t *testing.T) {
	ts := time.Now()
	period := OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
	}
	cs, err := period.CurrentStart()
	assert.True(t, cs.Equal(ts.AddDate(0, 0, -1)), "Expected period to return current start")
	assert.Nil(t, err, "Expected no error when getting current start")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 2)},
	}

	cs, err = period.CurrentStart()
	assert.Nil(t, cs, "Expected future period to not return current start")
	assert.NotNil(t, err, "Expected error when getting current start from future period")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -2)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
	}

	cs, err = period.CurrentStart()
	assert.Nil(t, cs, "Expected past period to not return current start")
	assert.NotNil(t, err, "Expected error when getting current start from past period")
}

func TestOncePeriodCurrentEnd(t *testing.T) {
	ts := time.Now()
	period := OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
	}
	ce, err := period.CurrentEnd()
	assert.True(t, ce.Equal(ts.AddDate(0, 0, 1)), "Expected period to return current end")
	assert.Nil(t, err, "Expected no error when getting current end")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 2)},
	}

	ce, err = period.CurrentEnd()
	assert.Nil(t, ce, "Expected future period to not return current end")
	assert.NotNil(t, err, "Expected error when getting current end from future period")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -2)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
	}

	ce, err = period.CurrentEnd()
	assert.Nil(t, ce, "Expected past period to not return current end")
	assert.NotNil(t, err, "Expected error when getting current end from past period")
}

func TestOncePeriodNextStart(t *testing.T) {
	ts := time.Now()
	period := OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
	}
	cs, err := period.NextStart()
	assert.Nil(t, cs, "Expected period to not return next start")
	assert.NotNil(t, err, "Expected error when getting next start")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 2)},
	}

	cs, err = period.NextStart()
	assert.True(t, cs.Equal(ts.AddDate(0, 0, 1)), "Expected future period to return next start")
	assert.Nil(t, err, "Expected no error when getting next start from future period")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -2)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
	}

	cs, err = period.NextStart()
	assert.Nil(t, cs, "Expected past period to not return next start")
	assert.NotNil(t, err, "Expected error when getting next start from past period")
}

func TestOncePeriodNextEnd(t *testing.T) {
	ts := time.Now()
	period := OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
	}
	cs, err := period.NextEnd()
	assert.Nil(t, cs, "Expected period to not return next end")
	assert.NotNil(t, err, "Expected error when getting next end")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 2)},
	}

	cs, err = period.NextEnd()
	assert.True(t, cs.Equal(ts.AddDate(0, 0, 2)), "Expected future period to return next end")
	assert.Nil(t, err, "Expected no error when getting next end from future period")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -2)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
	}

	cs, err = period.NextEnd()
	assert.Nil(t, cs, "Expected past period to not return next end")
	assert.NotNil(t, err, "Expected error when getting next end from past period")
}

func TestOncePeriodPreviousStart(t *testing.T) {
	ts := time.Now()
	period := OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
	}
	cs, err := period.PreviousStart()
	assert.Nil(t, cs, "Expected period to not return previous start")
	assert.NotNil(t, err, "Expected error when getting previous start")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 2)},
	}

	cs, err = period.PreviousStart()
	assert.Nil(t, cs, "Expected future period to not return previous start")
	assert.NotNil(t, err, "Expected error when getting previous start from future period")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -2)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
	}

	cs, err = period.PreviousStart()
	assert.True(t, cs.Equal(ts.AddDate(0, 0, -2)), "Expected past period to return previous start")
	assert.Nil(t, err, "Expected no error when getting previous start from past period")
}

func TestOncePeriodPreviousEnd(t *testing.T) {
	ts := time.Now()
	period := OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
	}
	cs, err := period.PreviousEnd()
	assert.Nil(t, cs, "Expected period to not return previous end")
	assert.NotNil(t, err, "Expected error when getting previous end")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, 1)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, 2)},
	}

	cs, err = period.PreviousEnd()
	assert.Nil(t, cs, "Expected future period to not return previous end")
	assert.NotNil(t, err, "Expected error when getting previous end from future period")

	period = OncePeriod{
		From: TimestampEdge{Timestamp: ts.AddDate(0, 0, -2)},
		To:   TimestampEdge{Timestamp: ts.AddDate(0, 0, -1)},
	}

	cs, err = period.PreviousEnd()
	assert.True(t, cs.Equal(ts.AddDate(0, 0, -1)), "Expected past period to return previous end")
	assert.Nil(t, err, "Expected no error when getting previous end from past period")
}

func TestTimestampEdgeBefore(t *testing.T) {
	layout := "2006-01-02 15:04"
	ts, _ := time.Parse(layout, "2025-04-28 12:00")
	edge := TimestampEdge{Timestamp: ts}
	past, _ := time.Parse(layout, "2025-04-21 07:35")
	future, _ := time.Parse(layout, "2025-04-29 12:00")

	assert.False(t, edge.Before(past), "Expected edge to not be before the past timestamp")
	assert.False(t, edge.Before(ts), "Expected edge to not be before the same timestamp")
	assert.True(t, edge.Before(future), "Expected edge to be before the future timestamp")
}

func TestTimestampEdgeAfter(t *testing.T) {
	layout := "2006-01-02 15:04"
	ts, _ := time.Parse(layout, "2025-04-28 12:00")
	edge := TimestampEdge{Timestamp: ts}
	past, _ := time.Parse(layout, "2025-04-21 07:35")
	future, _ := time.Parse(layout, "2025-04-29 12:00")

	assert.True(t, edge.After(past), "Expected edge to be after the past timestamp")
	assert.False(t, edge.After(ts), "Expected edge to not be after the same timestamp")
	assert.False(t, edge.After(future), "Expected edge to not be after the future timestamp")
}

func TestTimestampEdgeEqual(t *testing.T) {
	layout := "2006-01-02 15:04"
	ts, _ := time.Parse(layout, "2025-04-28 12:00")
	edge := TimestampEdge{Timestamp: ts}
	past, _ := time.Parse(layout, "2025-04-21 07:35")
	future, _ := time.Parse(layout, "2025-04-29 12:00")

	assert.False(t, edge.Equal(past), "Expected edge not to be equal the past timestamp")
	assert.True(t, edge.Equal(ts), "Expected edge to be equal the same timestamp")
	assert.False(t, edge.Equal(future), "Expected edge to not be equal the future timestamp")
}

func TestTimestampEdgeBeforeOrEqual(t *testing.T) {
	layout := "2006-01-02 15:04"
	ts, _ := time.Parse(layout, "2025-04-28 12:00")
	edge := TimestampEdge{Timestamp: ts}
	past, _ := time.Parse(layout, "2025-04-21 07:35")
	future, _ := time.Parse(layout, "2025-04-29 12:00")

	assert.False(t, edge.BeforeOrEqual(past), "Expected edge to not be before or equal the past timestamp")
	assert.True(t, edge.BeforeOrEqual(ts), "Expected edge to not be before or equal the same timestamp")
	assert.True(t, edge.BeforeOrEqual(future), "Expected edge to be before or equal the future timestamp")
}

func TestTimestampEdgeAfterOrEqual(t *testing.T) {
	layout := "2006-01-02 15:04"
	ts, _ := time.Parse(layout, "2025-04-28 12:00")
	edge := TimestampEdge{Timestamp: ts}
	past, _ := time.Parse(layout, "2025-04-21 07:35")
	future, _ := time.Parse(layout, "2025-04-29 12:00")

	assert.True(t, edge.AfterOrEqual(past), "Expected edge to be after or equal the past timestamp")
	assert.True(t, edge.AfterOrEqual(ts), "Expected edge to be after or equal the same timestamp")
	assert.False(t, edge.AfterOrEqual(future), "Expected edge to not be after or equal the future timestamp")
}
