package casoncelli

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	exampleJson := `{
   "periods":[
      {
         "name":"scheduled maintainance",
         "description":"update indexes",
         "type":"weekly",
         "from":{
            "day":"saturday",
            "hour":"23:00"
         },
         "to":{
            "day":"sunday",
            "hour":"07:00"
         }
      },
      {
         "name":"service interruption",
         "description":"as defined in mail 18/02/2025",
         "type":"once",
         "from":{
            "timestamp":"2025-02-20 12:30:00"
         },
         "to":{
            "timestamp":"2025-02-20 14:30:00"
         }
      }
   ]
}`
	var casoncelli Casoncelli
	err := json.Unmarshal([]byte(exampleJson), &casoncelli)
	assert.NoError(t, err, "Expected no error during unmarshalling")
	assert.Equal(t, 2, len(casoncelli.Periods), "Expected 2 periods to be unmarshalled")

	exp := Casoncelli{
		Periods: []Period{
			WeeklyPeriod{
				PeriodLabel: PeriodLabel{
					Name:        "scheduled maintainance",
					Description: "update indexes",
				},
				From: DayTimeEdge{
					Day:  time.Saturday,
					Hour: "23:00",
				},
				To: DayTimeEdge{
					Day:  time.Sunday,
					Hour: "07:00",
				},
			},
			OncePeriod{
				PeriodLabel: PeriodLabel{
					Name:        "service interruption",
					Description: "as defined in mail 18/02/2025",
				},
				From: TimestampEdge{
					Timestamp: time.Date(2025, 2, 20, 12, 30, 0, 0, time.Now().Location()),
				},
				To: TimestampEdge{
					Timestamp: time.Date(2025, 2, 20, 14, 30, 0, 0, time.Now().Location()),
				},
			},
		},
	}
	assert.True(t, casoncelli.Periods[0] == exp.Periods[0], "Expected result to contain the weekly period")
	assert.True(t, casoncelli.Periods[1] == exp.Periods[1], "Expected result to contain the once period")
}

type MockPeriod struct {
	result bool
}

func (m MockPeriod) Contains(t time.Time) bool {
	return m.result
}

func (m MockPeriod) ContainsNow() bool {
	return m.result
}

func (m MockPeriod) CurrentStart() (*time.Time, error) {
	return nil, nil
}

func (m MockPeriod) CurrentEnd() (*time.Time, error) {
	return nil, nil
}

func (m MockPeriod) NextStart() (*time.Time, error) {
	return nil, nil
}

func (m MockPeriod) NextEnd() (*time.Time, error) {
	return nil, nil
}

func (m MockPeriod) PreviousStart() (*time.Time, error) {
	return nil, nil
}

func (m MockPeriod) PreviousEnd() (*time.Time, error) {
	return nil, nil
}

func TestContains(t *testing.T) {
	c1 := Casoncelli{
		Periods: []Period{
			MockPeriod{result: true},
			MockPeriod{result: true},
		},
	}

	assert.True(t, c1.Contains(time.Now()), "Expected Contains to return true for c1")

	c2 := Casoncelli{
		Periods: []Period{
			MockPeriod{result: false},
			MockPeriod{result: false},
		},
	}
	assert.False(t, c2.Contains(time.Now()), "Expected Contains to return false for c2")

	c3 := Casoncelli{
		Periods: []Period{
			MockPeriod{result: false},
			MockPeriod{result: true},
		},
	}
	assert.True(t, c3.Contains(time.Now()), "Expected Contains to return true for c3")

	c4 := Casoncelli{
		Periods: []Period{
			MockPeriod{result: true},
			MockPeriod{result: false},
		},
	}
	assert.True(t, c4.Contains(time.Now()), "Expected Contains to return true for c4")

	c5 := Casoncelli{
		Periods: []Period{},
	}
	assert.False(t, c5.Contains(time.Now()), "Expected Contains to return false for c5")

	c6 := Casoncelli{
		Periods: nil,
	}
	assert.False(t, c6.Contains(time.Now()), "Expected Contains to return false for c6")
}

func TestContainsNow(t *testing.T) {
	c1 := Casoncelli{
		Periods: []Period{
			MockPeriod{result: true},
			MockPeriod{result: true},
		},
	}

	assert.True(t, c1.ContainsNow(), "Expected ContainsNow to return true for c1")

	c2 := Casoncelli{
		Periods: []Period{
			MockPeriod{result: false},
			MockPeriod{result: false},
		},
	}
	assert.False(t, c2.ContainsNow(), "Expected ContainsNow to return false for c2")

	c3 := Casoncelli{
		Periods: []Period{
			MockPeriod{result: false},
			MockPeriod{result: true},
		},
	}
	assert.True(t, c3.ContainsNow(), "Expected ContainsNow to return true for c3")

	c4 := Casoncelli{
		Periods: []Period{
			MockPeriod{result: true},
			MockPeriod{result: false},
		},
	}
	assert.True(t, c4.ContainsNow(), "Expected ContainsNow to return true for c4")

	c5 := Casoncelli{
		Periods: []Period{},
	}
	assert.False(t, c5.ContainsNow(), "Expected ContainsNow to return false for c5")

	c6 := Casoncelli{
		Periods: nil,
	}
	assert.False(t, c6.ContainsNow(), "Expected ContainsNow to return false for c6")
}
