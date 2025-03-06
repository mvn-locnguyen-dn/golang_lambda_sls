package nulls

import (
	"database/sql/driver"
	"encoding/json"
)

// Int32 null int32
type Int32 struct {
	Int32 int32
	Valid bool
}

// Value Marshal
func (i Int32) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return int64(i.Int32), nil
}

// Scan gorm scan
func (i *Int32) Scan(value interface{}) error {
	i.Valid, i.Int32 = value != nil, 0
	if i.Valid {
		i.Int32 = int32(value.(int64))
	}
	return nil
}

// MarshalJSON json MarshalJSON
func (i Int32) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Int32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON json UnmarshalJSON
func (i *Int32) UnmarshalJSON(data []byte) error {
	var v *int32
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	i.Valid = v != nil
	if i.Valid {
		i.Int32 = *v
	}

	return nil
}

// IsEqual compare between 2 Int32
func (i *Int32) IsEqual(j Int32) bool {
	return i.Valid == j.Valid && (!i.Valid || i.Int32 == j.Int32)
}
