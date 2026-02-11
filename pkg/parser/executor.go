package parser

import (
	"fmt"
	"friskadb/pkg/core"
	"strings"
)

// Executor executes parsed queries against a database
type Executor struct {
	db *core.Database
}

// NewExecutor creates a new executor instance
func NewExecutor(db *core.Database) *Executor {
	return &Executor{db: db}
}

// Execute runs a query and returns the result
func (e *Executor) Execute(query *Query) (interface{}, error) {
	switch query.Type {
	case QueryCreate:
		return e.executeCreate(query)
	case QuerySelect:
		return e.executeSelect(query)
	case QueryInsert:
		return e.executeInsert(query)
	case QueryUpdate:
		return e.executeUpdate(query)
	case QueryDelete:
		return e.executeDelete(query)
	case QueryDrop:
		return e.executeDrop(query)
	case QueryDescribe:
		return e.executeDescribe(query)
	case QueryShowTables:
		return e.executeShowTables(query)
	default:
		return nil, core.ErrInvalidQuery
	}
}

func (e *Executor) executeCreate(query *Query) (string, error) {
	err := e.db.CreateTable(query.TableName, query.Schema)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("✨ Table '%s' created successfully!", query.TableName), nil
}

func (e *Executor) executeSelect(query *Query) ([]core.Row, error) {
	table, err := e.db.GetTable(query.TableName)
	if err != nil {
		return nil, err
	}

	filter := e.buildFilter(query.Conditions)
	results, err := table.Select(query.Columns, filter)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (e *Executor) executeInsert(query *Query) (string, error) {
	table, err := e.db.GetTable(query.TableName)
	if err != nil {
		return "", err
	}

	// Build row from columns and values
	row := make(core.Row)
	for i, col := range query.Columns {
		if i < len(query.Values) {
			row[col] = query.Values[i]
		}
	}

	err = table.Insert(row)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("✅ Saved successfully! Total rows: %d", table.Count()), nil
}

func (e *Executor) executeUpdate(query *Query) (string, error) {
	table, err := e.db.GetTable(query.TableName)
	if err != nil {
		return "", err
	}

	filter := e.buildFilter(query.Conditions)
	count, err := table.Update(query.Updates, filter)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("✨ Updated %d row(s) successfully!", count), nil
}

func (e *Executor) executeDelete(query *Query) (string, error) {
	table, err := e.db.GetTable(query.TableName)
	if err != nil {
		return "", err
	}

	filter := e.buildFilter(query.Conditions)
	count, err := table.Delete(filter)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("🗑️ Deleted %d row(s). Remaining: %d", count, table.Count()), nil
}

func (e *Executor) executeDrop(query *Query) (string, error) {
	err := e.db.DropTable(query.TableName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("💥 Table '%s' dropped successfully!", query.TableName), nil
}

func (e *Executor) executeDescribe(query *Query) (*core.Table, error) {
	table, err := e.db.GetTable(query.TableName)
	if err != nil {
		return nil, err
	}

	return table, nil
}

func (e *Executor) executeShowTables(query *Query) ([]string, error) {
	tables := e.db.ListTables()
	return tables, nil
}

// buildFilter creates a filter function from conditions
func (e *Executor) buildFilter(cond *Condition) func(core.Row) bool {
	if cond == nil {
		return nil
	}

	return func(row core.Row) bool {
		return e.evaluateCondition(row, cond)
	}
}

func (e *Executor) evaluateCondition(row core.Row, cond *Condition) bool {
	if cond == nil {
		return true
	}

	value, exists := row[cond.Column]
	if !exists {
		return false
	}

	result := e.compareValues(value, cond.Operator, cond.Value)

	// Handle logical operators
	if cond.Next != nil {
		nextResult := e.evaluateCondition(row, cond.Next)
		if cond.Logic == "FRISAND" || cond.Logic == "AND" {
			return result && nextResult
		} else if cond.Logic == "FRISOR" || cond.Logic == "OR" {
			return result || nextResult
		}
	}

	return result
}

func (e *Executor) compareValues(left interface{}, op string, right interface{}) bool {
	switch op {
	case "=":
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
	case "!=":
		return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right)
	case ">":
		return e.numericCompare(left, right) > 0
	case "<":
		return e.numericCompare(left, right) < 0
	case ">=":
		return e.numericCompare(left, right) >= 0
	case "<=":
		return e.numericCompare(left, right) <= 0
	case "LIKE":
		return e.likeCompare(fmt.Sprintf("%v", left), fmt.Sprintf("%v", right))
	case "NOTNULL":
		return left != nil
	case "NULL":
		return left == nil
	default:
		return false
	}
}

func (e *Executor) numericCompare(left, right interface{}) int {
	l := e.toFloat(left)
	r := e.toFloat(right)

	if l > r {
		return 1
	} else if l < r {
		return -1
	}
	return 0
}

func (e *Executor) toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

func (e *Executor) likeCompare(str, pattern string) bool {
	// Simple LIKE implementation with % wildcard
	pattern = strings.ReplaceAll(pattern, "%", ".*")
	pattern = "^" + pattern + "$"

	// Simple contains check for now
	if strings.HasPrefix(pattern, "^.*") && strings.HasSuffix(pattern, ".*$") {
		// %pattern% - contains
		substr := strings.TrimPrefix(pattern, "^.*")
		substr = strings.TrimSuffix(substr, ".*$")
		return strings.Contains(str, substr)
	} else if strings.HasPrefix(pattern, "^.*") {
		// %pattern - ends with
		suffix := strings.TrimPrefix(pattern, "^.*")
		suffix = strings.TrimSuffix(suffix, "$")
		return strings.HasSuffix(str, suffix)
	} else if strings.HasSuffix(pattern, ".*$") {
		// pattern% - starts with
		prefix := strings.TrimPrefix(pattern, "^")
		prefix = strings.TrimSuffix(prefix, ".*$")
		return strings.HasPrefix(str, prefix)
	}

	return str == strings.Trim(pattern, "^$")
}
