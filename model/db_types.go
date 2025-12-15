package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// StringSlice is a custom type for []string that can be scanned from JSON in database
type StringSlice []string

// Scan implements sql.Scanner interface
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported type for StringSlice: %T", value)
	}

	// Handle empty or null JSON
	if len(bytes) == 0 || string(bytes) == "null" {
		*s = []string{}
		return nil
	}

	return json.Unmarshal(bytes, s)
}

// Value implements driver.Valuer interface
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}
