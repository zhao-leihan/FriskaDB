package core

import (
	"fmt"
	"sync"
)

// Database represents the main FriskaDB database
type Database struct {
	Name   string            `json:"name"`
	Tables map[string]*Table `json:"tables"`
	mu     sync.RWMutex
}

// NewDatabase creates a new database instance
func NewDatabase(name string) *Database {
	return &Database{
		Name:   name,
		Tables: make(map[string]*Table),
	}
}

// CreateTable creates a new table in the database
func (db *Database) CreateTable(name string, columns []Column) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.Tables[name]; exists {
		return fmt.Errorf("%w: %s", ErrTableExists, name)
	}

	db.Tables[name] = NewTable(name, columns)
	return nil
}

// GetTable retrieves a table by name
func (db *Database) GetTable(name string) (*Table, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	table, exists := db.Tables[name]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrTableNotFound, name)
	}

	return table, nil
}

// DropTable removes a table from the database
func (db *Database) DropTable(name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.Tables[name]; !exists {
		return fmt.Errorf("%w: %s", ErrTableNotFound, name)
	}

	delete(db.Tables, name)
	return nil
}

// ListTables returns all table names
func (db *Database) ListTables() []string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	names := make([]string, 0, len(db.Tables))
	for name := range db.Tables {
		names = append(names, name)
	}
	return names
}

// TableExists checks if a table exists
func (db *Database) TableExists(name string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()

	_, exists := db.Tables[name]
	return exists
}
