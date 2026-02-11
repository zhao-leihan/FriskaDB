package core

import (
	"fmt"
	"sync"
)

// Table represents a database table
type Table struct {
	Name   string `json:"name"`
	Schema Schema `json:"schema"`
	Rows   []Row  `json:"rows"`
	mu     sync.RWMutex
}

// NewTable creates a new table with the given schema
func NewTable(name string, columns []Column) *Table {
	return &Table{
		Name: name,
		Schema: Schema{
			Columns: columns,
		},
		Rows: make([]Row, 0),
	}
}

// Insert adds a new row to the table
func (t *Table) Insert(row Row) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Validate row against schema
	if err := t.validateRow(row); err != nil {
		return err
	}

	t.Rows = append(t.Rows, row)
	return nil
}

// Select returns rows matching the filter function
func (t *Table) Select(columns []string, filter func(Row) bool) ([]Row, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var results []Row

	for _, row := range t.Rows {
		if filter == nil || filter(row) {
			// If specific columns requested, return only those
			if len(columns) > 0 && columns[0] != "*" {
				filteredRow := make(Row)
				for _, col := range columns {
					if val, exists := row[col]; exists {
						filteredRow[col] = val
					} else {
						return nil, fmt.Errorf("%w: %s", ErrColumnNotFound, col)
					}
				}
				results = append(results, filteredRow)
			} else {
				// Return all columns
				results = append(results, row)
			}
		}
	}

	return results, nil
}

// Update modifies rows matching the filter
func (t *Table) Update(updates Row, filter func(Row) bool) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	count := 0
	for i := range t.Rows {
		if filter == nil || filter(t.Rows[i]) {
			// Apply updates
			for key, value := range updates {
				t.Rows[i][key] = value
			}
			count++
		}
	}

	return count, nil
}

// Delete removes rows matching the filter
func (t *Table) Delete(filter func(Row) bool) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var remaining []Row
	count := 0

	for _, row := range t.Rows {
		if filter == nil || filter(row) {
			count++
		} else {
			remaining = append(remaining, row)
		}
	}

	t.Rows = remaining
	return count, nil
}

// Count returns the number of rows in the table
func (t *Table) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.Rows)
}

// validateRow checks if row matches table schema
func (t *Table) validateRow(row Row) error {
	// Check if all required columns are present
	columnMap := make(map[string]DataType)
	for _, col := range t.Schema.Columns {
		columnMap[col.Name] = col.Type
	}

	for colName := range row {
		if _, exists := columnMap[colName]; !exists {
			return fmt.Errorf("%w: %s", ErrColumnNotFound, colName)
		}
	}

	return nil
}

// GetColumnNames returns all column names
func (t *Table) GetColumnNames() []string {
	names := make([]string, len(t.Schema.Columns))
	for i, col := range t.Schema.Columns {
		names[i] = col.Name
	}
	return names
}
