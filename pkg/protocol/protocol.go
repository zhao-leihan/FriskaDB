package protocol

import (
	"encoding/json"
	"RayhanDB/pkg/core"
)

// Request represents a client request to the server
type Request struct {
	ID    string `json:"id"`    // Unique request ID for tracking
	Query string `json:"query"` // Friska query string
	Auth  Auth   `json:"auth"`  // Authentication credentials
}

// Auth holds authentication credentials
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Response represents server response to client
type Response struct {
	ID      string      `json:"id"`      // Matching request ID
	Success bool        `json:"success"` // Query execution status
	Data    interface{} `json:"data"`    // Query results
	Error   string      `json:"error"`   // Error message if failed
	Message string      `json:"message"` // Success message
}

// RowData represents a single row for serialization
type RowData map[string]interface{}

// TableInfo represents table schema information
type TableInfo struct {
	Name     string       `json:"name"`
	Columns  []ColumnInfo `json:"columns"`
	RowCount int          `json:"row_count"`
}

// ColumnInfo represents column metadata
type ColumnInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// EncodeRequest serializes a request to JSON
func EncodeRequest(req *Request) ([]byte, error) {
	return json.Marshal(req)
}

// DecodeRequest deserializes JSON to request
func DecodeRequest(data []byte) (*Request, error) {
	var req Request
	err := json.Unmarshal(data, &req)
	return &req, err
}

// EncodeResponse serializes a response to JSON
func EncodeResponse(resp *Response) ([]byte, error) {
	return json.Marshal(resp)
}

// DecodeResponse deserializes JSON to response
func DecodeResponse(data []byte) (*Response, error) {
	var resp Response
	err := json.Unmarshal(data, &resp)
	return &resp, err
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(requestID string, data interface{}, message string) *Response {
	return &Response{
		ID:      requestID,
		Success: true,
		Data:    data,
		Message: message,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(requestID string, err error) *Response {
	return &Response{
		ID:      requestID,
		Success: false,
		Error:   err.Error(),
	}
}

// ConvertRows converts core.Row slice to RowData slice for serialization
func ConvertRows(rows []core.Row) []RowData {
	result := make([]RowData, len(rows))
	for i, row := range rows {
		result[i] = RowData(row)
	}
	return result
}

// ConvertTableInfo converts core.Table to TableInfo
func ConvertTableInfo(table *core.Table) TableInfo {
	columns := make([]ColumnInfo, len(table.Schema.Columns))
	for i, col := range table.Schema.Columns {
		columns[i] = ColumnInfo{
			Name: col.Name,
			Type: string(col.Type),
		}
	}

	return TableInfo{
		Name:     table.Name,
		Columns:  columns,
		RowCount: table.Count(),
	}
}
