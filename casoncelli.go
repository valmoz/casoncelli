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
		//Timezone *string           `json:"timezone,omitempty"`
	}

	var rawObj rawCasoncelli
	if err := json.Unmarshal(data, &rawObj); err != nil {
		return err
	}

	periods := []Period{}
	for _, raw := range rawObj.Periods {
		var peek struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &peek); err != nil {
			return err
		}

		switch peek.Type {
		case "weekly":
			var wp WeeklyPeriod
			if err := json.Unmarshal(raw, &wp); err != nil {
				return err
			}
			periods = append(periods, wp)
		case "once":
			var op OncePeriod
			if err := json.Unmarshal(raw, &op); err != nil {
				return err
			}
			periods = append(periods, op)
		default:
			return fmt.Errorf("unknown period type: %s", peek.Type)
		}
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
