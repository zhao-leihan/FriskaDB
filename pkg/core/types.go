package core

import "time"

// DataType represents supported column types in RayhanDB
type DataType string

const (
	TypeText    DataType = "TEXT"
	TypeNumber  DataType = "NUMBER"
	TypeBoolean DataType = "BOOLEAN"
	TypeDate    DataType = "DATE"
)

// Column represents a table column definition
type Column struct {
	Name string   `json:"name"`
	Type DataType `json:"type"`
}

// Row represents a single row of data
type Row map[string]interface{}

// Schema represents a table schema
type Schema struct {
	Columns []Column `json:"columns"`
}

// Value represents a typed value
type Value struct {
	Raw  interface{}
	Type DataType
}

// ParseValue converts raw value to typed Value
func ParseValue(raw interface{}, dataType DataType) (*Value, error) {
	switch dataType {
	case TypeText:
		if str, ok := raw.(string); ok {
			return &Value{Raw: str, Type: TypeText}, nil
		}
		return nil, ErrInvalidType
	case TypeNumber:
		switch v := raw.(type) {
		case float64, int, int64:
			return &Value{Raw: v, Type: TypeNumber}, nil
		}
		return nil, ErrInvalidType
	case TypeBoolean:
		if b, ok := raw.(bool); ok {
			return &Value{Raw: b, Type: TypeBoolean}, nil
		}
		return nil, ErrInvalidType
	case TypeDate:
		switch v := raw.(type) {
		case time.Time:
			return &Value{Raw: v, Type: TypeDate}, nil
		case string:
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, err
			}
			return &Value{Raw: t, Type: TypeDate}, nil
		}
		return nil, ErrInvalidType
	}
	return nil, ErrInvalidType
}
