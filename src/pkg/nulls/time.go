package nulls

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Time null time
type Time struct {
	Time  time.Time
	Valid bool
}

// Value Marshal
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// Scan gorm scan
func (t *Time) Scan(value interface{}) error {
	t.Time, t.Valid = value.(time.Time)
	return nil
}

// MarshalJSON json MarshalJSON
func (t Time) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.Time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON json UnmarshalJSON
func (t *Time) UnmarshalJSON(data []byte) error {
	var v *time.Time
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	t.Valid = false
	if v != nil {
		t.Valid = true
		t.Time = *v
	}

	return nil
}
