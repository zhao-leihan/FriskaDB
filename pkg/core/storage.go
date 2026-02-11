package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Storage handles database persistence
type Storage struct {
	DataDir string
}

// NewStorage creates a new storage instance
func NewStorage(dataDir string) *Storage {
	return &Storage{
		DataDir: dataDir,
	}
}

// Save persists a database to disk
func (s *Storage) Save(db *Database) error {
	// Ensure data directory exists
	if err := os.MkdirAll(s.DataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	filename := filepath.Join(s.DataDir, db.Name+".fris")

	// Marshal database to JSON
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal database: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write database file: %w", err)
	}

	return nil
}

// Load reads a database from disk
func (s *Storage) Load(name string) (*Database, error) {
	filename := filepath.Join(s.DataDir, name+".fris")

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, ErrDatabaseNotFound
	}

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read database file: %w", err)
	}

	// Unmarshal JSON
	var db Database
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, fmt.Errorf("failed to unmarshal database: %w", err)
	}

	return &db, nil
}

// Delete removes a database file from disk
func (s *Storage) Delete(name string) error {
	filename := filepath.Join(s.DataDir, name+".fris")

	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return ErrDatabaseNotFound
		}
		return fmt.Errorf("failed to delete database file: %w", err)
	}

	return nil
}

// Exists checks if a database file exists
func (s *Storage) Exists(name string) bool {
	filename := filepath.Join(s.DataDir, name+".fris")
	_, err := os.Stat(filename)
	return err == nil
}
