# Casoncelli

Casoncelli is a simple Go Library for managing recurring time periods.

The library takes its name from a type of filled pasta (called "[Casoncelli](https://en.wikipedia.org/wiki/Casoncelli)") typical of northern Italy. As the filling can be inside or outside a single Casoncello piece, the library can be used to check if a timestamp is contained in or not in a single Period.

Its main objective is to manage automatically scheduled maintenance periods, when the application must be unavailable or a certain operation should not be triggered.

## Functionalities

Casoncelli can currently manage two different types of periods:

- **Weekly Periods**: Periods that repeat themselves every 7 days
- **Once Periods**: These are single events that happen only once and never again

Each period is defined by its **Edges**, which declare when a period starts or finishes.

The library makes it possible to combine multiple periods; a timestamp will be considered contained if at least one of the periods contains it.

The periods can be declared directly from the code or by a JSON string; this makes it possible to store the configuration somewhere and load it dynamically when needed.

### Weekly Periods

A Weekly Period is defined by day/hour edges, for example:

```json
{
  "name": "scheduled maintenance",
  "description": "update indexes",
  "type": "weekly",
  "from": {
    "day": "saturday",
    "hour": "23:00"
  },
  "to": {
    "day": "sunday",
    "hour": "07:00"
  }
}
```

In this case, the period starts every Saturday at 23:00 and ends the next Sunday at 07:00 local time.

### Once Periods

A Once period is defined by timestamp edges, for example:

```json
{
  "name": "service interruption",
  "description": "as defined in mail 18/02/2025",
  "type": "once",
  "from": {
    "timestamp": "2025-02-20 12:30:00"
  },
  "to": {
    "timestamp": "2025-02-20 14:30:00"
  }
}
```

In this case, the period starts at 12:30 on 20 February 2025 and ends at 14:30 on the same day.

## Installation

```bash
go get github.com/valmoz/casoncelli
```

## Usage

### Base example

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/valmoz/casoncelli"
)

func main() {
    // Period definition
    jsonConfig := `[
        {
            "name": "scheduled maintenance",
            "description": "update indexes",
            "type": "weekly",
            "from": {
                "day": "saturday",
                "hour": "23:00"
            },
            "to": {
                "day": "sunday",
                "hour": "07:00"
            }
        },
        {
            "name": "service interruption",
            "description": "as defined in mail 18/02/2025",
            "type": "once",
            "from": {
                "timestamp": "2025-02-20 12:30:00"
            },
            "to": {
                "timestamp": "2025-02-20 14:30:00"
            }
        }
    ]`

    var casoncelli casoncelli.Casoncelli
    err := casoncelli.UnmarshalJSON([]byte(jsonConfig), &casoncelli)
    if err != nil {
        panic(err)
    }

    // Checking if the system is currently offline
    if casoncelli.ContainsNow() {
        fmt.Println("The service is currently offline")
    } else {
        fmt.Println("The service is currently online")
    }

    // Checking a specific time
    specificTime := time.Date(2025, 2, 20, 13, 0, 0, 0, time.UTC)
    if casoncelli.Contains(specificTime) {
        fmt.Println("The specified time is contained in the period")
    }
}
```

## API

### `Casoncelli` methods

- `Contains(t time.Time) bool`: Returns true if `t` is included in at least one of the periods
- `ContainsNow() bool`: Returns true if the current moment is included in the periods

### `Period` methods

- `Contains(t time.Time) bool`: Returns true if the moment `t` is in the period
- `ContainsNow() bool`: Returns true if the current moment `t` is in the period
- `CurrentStart() (*time.Time, error)`: If the period is currently active, returns the start of the period
- `CurrentEnd() (*time.Time, error)`: If the period is currently active, returns the end of the period
- `NextStart() (*time.Time, error)`: Returns the start of the next period
- `NextEnd() (*time.Time, error)`: Returns the end of the next period
- `PreviousStart() (*time.Time, error)`: Returns the start of the previous period
- `PreviousEnd() (*time.Time, error)`: Returns the end of the previous period

## Test

The tests can be executed with:

```bash
go test ./...
```

## License

This library is distributed under the MIT license found in the [LICENSE](./LICENSE)
file.
