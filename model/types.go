package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

type NullUTC8Time struct {
	Time  time.Time
	Valid bool
}

func (n *NullUTC8Time) Scan(value interface{}) error {
	if value == nil {
		n.Time, n.Valid = time.Now(), false
		return nil
	}
	n.Valid = true
	if s, ok := value.(time.Time); ok {
		n.Time = s
	} else {
		return errors.New("NonTimeData")
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullUTC8Time) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}
