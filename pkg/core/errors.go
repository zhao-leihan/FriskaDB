package core

import "errors"

var (
	// Database errors
	ErrDatabaseNotFound = errors.New("😅 Oops! Database doesn't exist yet")
	ErrDatabaseExists   = errors.New("⚠️ Database already exists!")

	// Table errors
	ErrTableNotFound = errors.New("😅 Oops! Table doesn't exist yet")
	ErrTableExists   = errors.New("⚠️ Table already exists!")
	ErrInvalidSchema = errors.New("❌ Invalid table schema")

	// Data errors
	ErrInvalidType    = errors.New("❌ Invalid data type")
	ErrColumnNotFound = errors.New("😅 Column doesn't exist")
	ErrMissingColumn  = errors.New("❌ Required column is missing")
	ErrRowNotFound    = errors.New("😅 No rows found")

	// Query errors
	ErrInvalidQuery = errors.New("❌ Invalid query syntax")
	ErrEmptyQuery   = errors.New("😅 Query is empty")
)
