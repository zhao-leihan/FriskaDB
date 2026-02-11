package parser

import (
	"fmt"
	"friskadb/pkg/core"
	"strconv"
	"strings"
)

// QueryType represents the type of query
type QueryType int

const (
	QueryUnknown QueryType = iota
	QueryCreate
	QuerySelect
	QueryInsert
	QueryUpdate
	QueryDelete
	QueryDrop
	QueryDescribe
	QueryShowTables
)

// Query represents a parsed SQL query
type Query struct {
	Type       QueryType
	TableName  string
	Columns    []string
	Values     []interface{}
	Conditions *Condition
	Updates    map[string]interface{}
	Schema     []core.Column
}

// Condition represents a WHERE clause condition
type Condition struct {
	Column   string
	Operator string
	Value    interface{}
	Logic    string // AND, OR
	Next     *Condition
}

// Parser parses Friska queries
type Parser struct {
	lexer     *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

// NewParser creates a new parser instance
func NewParser(input string) *Parser {
	p := &Parser{
		lexer:  NewLexer(input),
		errors: []string{},
	}

	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// Parse parses the input and returns a Query
func (p *Parser) Parse() (*Query, error) {
	var query *Query
	var err error

	switch p.curToken.Type {
	case TokenFrisrate:
		query, err = p.parseCreate()
	case TokenFrislect:
		query, err = p.parseSelect()
	case TokenFrisert:
		query, err = p.parseInsert()
	case TokenFrisdate:
		query, err = p.parseUpdate()
	case TokenFrislete:
		query, err = p.parseDelete()
	case TokenFrisdrop:
		query, err = p.parseDrop()
	case TokenFrisc:
		query, err = p.parseDescribe()
	case TokenFrisshow:
		query, err = p.parseShowTables()
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.curToken.Literal)
	}

	if err != nil {
		return nil, err
	}

	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parse errors: %s", strings.Join(p.errors, "; "))
	}

	return query, nil
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors, fmt.Sprintf("expected %v, got %v", t, p.peekToken.Type))
	return false
}

// parseCreate parses FRISRATE FRISKABLE table_name (columns...)
func (p *Parser) parseCreate() (*Query, error) {
	query := &Query{Type: QueryCreate}

	// Expect FRISKABLE
	if !p.expectPeek(TokenFriskable) {
		return nil, fmt.Errorf("expected FRISKABLE after FRISRATE")
	}

	// Expect table name
	if !p.expectPeek(TokenIdent) {
		return nil, fmt.Errorf("expected table name")
	}
	query.TableName = p.curToken.Literal

	// Expect (
	if !p.expectPeek(TokenLeftParen) {
		return nil, fmt.Errorf("expected ( after table name")
	}

	// Parse columns
	query.Schema = p.parseColumnDefinitions()

	// Expect )
	if !p.expectPeek(TokenRightParen) {
		return nil, fmt.Errorf("expected ) after column definitions")
	}

	return query, nil
}

func (p *Parser) parseColumnDefinitions() []core.Column {
	columns := []core.Column{}

	for p.peekToken.Type != TokenRightParen && p.peekToken.Type != TokenEOF {
		p.nextToken()

		colName := p.curToken.Literal

		// Expect type
		p.nextToken()
		colType := core.DataType(strings.ToUpper(p.curToken.Literal))

		columns = append(columns, core.Column{
			Name: colName,
			Type: colType,
		})

		// Check for comma
		if p.peekToken.Type == TokenComma {
			p.nextToken()
		}
	}

	return columns
}

// parseSelect parses FRISLECT columns FRISFROM table [FRISWHERE condition]
func (p *Parser) parseSelect() (*Query, error) {
	query := &Query{Type: QuerySelect}

	// Parse columns
	query.Columns = p.parseColumnList()

	// Expect FRISFROM
	if !p.expectPeek(TokenFrisfrom) {
		return nil, fmt.Errorf("expected FRISFROM")
	}

	// Expect table name
	if !p.expectPeek(TokenIdent) {
		return nil, fmt.Errorf("expected table name")
	}
	query.TableName = p.curToken.Literal

	// Optional WHERE clause
	if p.peekToken.Type == TokenFriswhere {
		p.nextToken()
		cond, err := p.parseCondition()
		if err != nil {
			return nil, err
		}
		query.Conditions = cond
	}

	return query, nil
}

func (p *Parser) parseColumnList() []string {
	columns := []string{}

	for {
		p.nextToken()

		if p.curToken.Type == TokenAsterisk {
			return []string{"*"}
		}

		if p.curToken.Type == TokenIdent {
			columns = append(columns, p.curToken.Literal)
		}

		if p.peekToken.Type != TokenComma {
			break
		}
		p.nextToken() // consume comma
	}

	return columns
}

// parseInsert parses FRISERT FRISINTO table (columns) FRISVALUES (values)
func (p *Parser) parseInsert() (*Query, error) {
	query := &Query{Type: QueryInsert}

	// Expect FRISINTO
	if !p.expectPeek(TokenFrisinto) {
		return nil, fmt.Errorf("expected FRISINTO")
	}

	// Expect table name
	if !p.expectPeek(TokenIdent) {
		return nil, fmt.Errorf("expected table name")
	}
	query.TableName = p.curToken.Literal

	// Expect (
	if !p.expectPeek(TokenLeftParen) {
		return nil, fmt.Errorf("expected (")
	}

	// Parse columns
	query.Columns = p.parseIdentList()

	// Expect )
	if !p.expectPeek(TokenRightParen) {
		return nil, fmt.Errorf("expected )")
	}

	// Expect FRISVALUES
	if !p.expectPeek(TokenFrisvalues) {
		return nil, fmt.Errorf("expected FRISVALUES")
	}

	// Expect (
	if !p.expectPeek(TokenLeftParen) {
		return nil, fmt.Errorf("expected (")
	}

	// Parse values
	query.Values = p.parseValueList()

	// Expect )
	if !p.expectPeek(TokenRightParen) {
		return nil, fmt.Errorf("expected )")
	}

	return query, nil
}

func (p *Parser) parseIdentList() []string {
	idents := []string{}

	for p.peekToken.Type != TokenRightParen {
		p.nextToken()
		idents = append(idents, p.curToken.Literal)

		if p.peekToken.Type == TokenComma {
			p.nextToken()
		}
	}

	return idents
}

func (p *Parser) parseValueList() []interface{} {
	values := []interface{}{}

	for p.peekToken.Type != TokenRightParen {
		p.nextToken()

		val := p.parseValue()
		values = append(values, val)

		if p.peekToken.Type == TokenComma {
			p.nextToken()
		}
	}

	return values
}

func (p *Parser) parseValue() interface{} {
	switch p.curToken.Type {
	case TokenString:
		return p.curToken.Literal
	case TokenNumber:
		if num, err := strconv.ParseFloat(p.curToken.Literal, 64); err == nil {
			return num
		}
		return p.curToken.Literal
	case TokenBoolean:
		return strings.ToLower(p.curToken.Literal) == "true"
	default:
		return p.curToken.Literal
	}
}

// parseUpdate parses FRISDATE table FRISSET col=val [FRISWHERE condition]
func (p *Parser) parseUpdate() (*Query, error) {
	query := &Query{
		Type:    QueryUpdate,
		Updates: make(map[string]interface{}),
	}

	// Expect table name
	if !p.expectPeek(TokenIdent) {
		return nil, fmt.Errorf("expected table name")
	}
	query.TableName = p.curToken.Literal

	// Expect FRISSET
	if !p.expectPeek(TokenFrisset) {
		return nil, fmt.Errorf("expected FRISSET")
	}

	// Parse assignments
	for {
		p.nextToken()
		colName := p.curToken.Literal

		// Expect =
		if !p.expectPeek(TokenEquals) {
			return nil, fmt.Errorf("expected = after column name")
		}

		p.nextToken()
		value := p.parseValue()
		query.Updates[colName] = value

		if p.peekToken.Type != TokenComma {
			break
		}
		p.nextToken() // consume comma
	}

	// Optional WHERE
	if p.peekToken.Type == TokenFriswhere {
		p.nextToken()
		cond, err := p.parseCondition()
		if err != nil {
			return nil, err
		}
		query.Conditions = cond
	}

	return query, nil
}

// parseDelete parses FRISLETE FRISFROM table [FRISWHERE condition]
func (p *Parser) parseDelete() (*Query, error) {
	query := &Query{Type: QueryDelete}

	// Expect FRISFROM
	if !p.expectPeek(TokenFrisfrom) {
		return nil, fmt.Errorf("expected FRISFROM")
	}

	// Expect table name
	if !p.expectPeek(TokenIdent) {
		return nil, fmt.Errorf("expected table name")
	}
	query.TableName = p.curToken.Literal

	// Optional WHERE
	if p.peekToken.Type == TokenFriswhere {
		p.nextToken()
		cond, err := p.parseCondition()
		if err != nil {
			return nil, err
		}
		query.Conditions = cond
	}

	return query, nil
}

// parseDrop parses FRISDROP FRISKABLE table
func (p *Parser) parseDrop() (*Query, error) {
	query := &Query{Type: QueryDrop}

	// Expect FRISKABLE
	if !p.expectPeek(TokenFriskable) {
		return nil, fmt.Errorf("expected FRISKABLE")
	}

	// Expect table name
	if !p.expectPeek(TokenIdent) {
		return nil, fmt.Errorf("expected table name")
	}
	query.TableName = p.curToken.Literal

	return query, nil
}

// parseDescribe parses FRISC table
func (p *Parser) parseDescribe() (*Query, error) {
	query := &Query{Type: QueryDescribe}

	// Expect table name
	if !p.expectPeek(TokenIdent) {
		return nil, fmt.Errorf("expected table name")
	}
	query.TableName = p.curToken.Literal

	return query, nil
}

// parseShowTables parses FRISSHOW FRISKABLES
func (p *Parser) parseShowTables() (*Query, error) {
	query := &Query{Type: QueryShowTables}

	// Expect FRISKABLES
	if !p.expectPeek(TokenFriskables) {
		return nil, fmt.Errorf("expected FRISKABLES")
	}

	return query, nil
}

// parseCondition parses WHERE conditions
func (p *Parser) parseCondition() (*Condition, error) {
	cond := &Condition{}

	// Expect column name
	p.nextToken()
	cond.Column = p.curToken.Literal

	// Expect operator
	p.nextToken()
	cond.Operator = p.getOperator()

	// Expect value
	p.nextToken()
	cond.Value = p.parseValue()

	// Check for logical operator (AND/OR)
	if p.peekToken.Type == TokenFrisand || p.peekToken.Type == TokenFrisor {
		p.nextToken()
		cond.Logic = strings.ToUpper(p.curToken.Literal)

		// Parse next condition
		nextCond, err := p.parseCondition()
		if err != nil {
			return nil, err
		}
		cond.Next = nextCond
	}

	return cond, nil
}

func (p *Parser) getOperator() string {
	switch p.curToken.Type {
	case TokenEquals:
		return "="
	case TokenNotEquals:
		return "!="
	case TokenGreater, TokenAbove:
		return ">"
	case TokenLess, TokenBelow:
		return "<"
	case TokenGreaterEq, TokenAtleast:
		return ">="
	case TokenLessEq, TokenAtmost:
		return "<="
	case TokenFrislove:
		return "LIKE"
	case TokenFrisamong:
		return "IN"
	case TokenFrisxist:
		return "NOTNULL"
	case TokenFrismiss:
		return "NULL"
	default:
		return p.curToken.Literal
	}
}
