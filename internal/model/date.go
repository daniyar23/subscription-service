package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// MonthYear represents a month and year as a time.Time value.
type MonthYear struct {
	time.Time
}

const layout = "01-2006"

// UnmarshalJSON parses a JSON string into a MonthYear value.
func (m *MonthYear) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	m.Time = t
	return nil
}

// MarshalJSON serializes a MonthYear value to a JSON string.
func (m MonthYear) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Time.Format(layout))
}

// Value converts a MonthYear value to a driver.Value.
func (m MonthYear) Value() (driver.Value, error) {
	return m.Time, nil
}

// Scan converts a driver.Value to a MonthYear value.
func (m *MonthYear) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot scan MonthYear")
	}
	m.Time = t
	return nil
}
