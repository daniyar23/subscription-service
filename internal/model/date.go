package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type MonthYear struct {
	time.Time
}

const layout = "01-2006"

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

func (m MonthYear) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Time.Format(layout))
}

func (m MonthYear) Value() (driver.Value, error) {
	return m.Time, nil
}

func (m *MonthYear) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot scan MonthYear")
	}
	m.Time = t
	return nil
}
