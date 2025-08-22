package casoncelli

import (
	"encoding/json"
	"fmt"
	"time"
)

type Casoncelli struct {
	Periods []Period `json:"periods"`
	//Timezone *time.Location `json:"timezone,omitempty"`
}

func (c *Casoncelli) UnmarshalJSON(data []byte) error {
	type rawCasoncelli struct {
		Periods []json.RawMessage `json:"periods"`
	}

	var rawObj rawCasoncelli
	if err := json.Unmarshal(data, &rawObj); err != nil {
		return err
	}

	periodTypes := map[string]func(json.RawMessage) (Period, error){
		"weekly": unmarshalPeriod[WeeklyPeriod],
		"daily":  unmarshalPeriod[DailyPeriod],
		"once":   unmarshalPeriod[OncePeriod],
		"never":  unmarshalPeriod[NeverPeriod],
		"always": unmarshalPeriod[AlwaysPeriod],
	}

	periods := []Period{}
	for _, raw := range rawObj.Periods {
		var peek struct {
			Type string `json:"type"`
		}

		if err := json.Unmarshal(raw, &peek); err != nil {
			return err
		}

		unmarshaler, exists := periodTypes[peek.Type]
		if !exists {
			return fmt.Errorf("unknown period type: %s", peek.Type)
		}

		period, err := unmarshaler(raw)
		if err != nil {
			return err
		}

		periods = append(periods, period)
	}

	c.Periods = periods
	return nil
}

func (c *Casoncelli) Contains(t time.Time) bool {
	for _, period := range c.Periods {
		if period.Contains(t) {
			return true
		}
	}
	return false
}

func (c *Casoncelli) ContainsNow() bool {
	for _, period := range c.Periods {
		if period.ContainsNow() {
			return true
		}
	}
	return false
}

func unmarshalPeriod[T Period](raw json.RawMessage) (Period, error) {
	var period T
	if err := json.Unmarshal(raw, &period); err != nil {
		return nil, err
	}
	return period, nil
}
