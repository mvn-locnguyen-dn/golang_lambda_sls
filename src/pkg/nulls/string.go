package nulls

import (
	"database/sql/driver"
	"encoding/json"
)

// String null string
type String struct {
	String string
	Valid  bool
}

// Value Marshal
func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.String, nil
}

// Scan gorm scan
func (s *String) Scan(value interface{}) error {
	s.String, s.Valid = value.(string)
	return nil
}

// MarshalJSON json MarshalJSON
func (s String) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON json UnmarshalJSON
func (s *String) UnmarshalJSON(data []byte) error {
	var v *string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	s.Valid = v != nil
	if s.Valid {
		s.String = *v
	}

	return nil
}
