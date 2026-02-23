package parser

import (
	"RayhanDB/pkg/core"
	"fmt"
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
	case TokenRAYRATE:
		query, err = p.parseCreate()
	case TokenRAYLECT:
		query, err = p.parseSelect()
	case TokenRAYERT:
		query, err = p.parseInsert()
	case TokenRAYDATE:
		query, err = p.parseUpdate()
	case TokenRAYLETE:
		query, err = p.parseDelete()
	case TokenRAYDROP:
		query, err = p.parseDrop()
	case TokenRAYC:
		query, err = p.parseDescribe()
	case TokenRAYSHOW:
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

// parseCreate parses RAYRATE RAYTABLE table_name (columns...)
func (p *Parser) parseCreate() (*Query, error) {
	query := &Query{Type: QueryCreate}

	// Expect RAYTABLE
	if !p.expectPeek(TokenRAYTABLE) {
		return nil, fmt.Errorf("expected RAYTABLE after RAYRATE")
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

// parseSelect parses RAYLECT columns RAYFROM table [RAYWHERE condition]
func (p *Parser) parseSelect() (*Query, error) {
	query := &Query{Type: QuerySelect}

	// Parse columns
	query.Columns = p.parseColumnList()

	// Expect RAYFROM
	if !p.expectPeek(TokenRAYFROM) {
		return nil, fmt.Errorf("expected RAYFROM")
	}

	// Expect table name
	if !p.expectPeek(TokenIdent) {
		return nil, fmt.Errorf("expected table name")
	}
	query.TableName = p.curToken.Literal

	// Optional WHERE clause
	if p.peekToken.Type == TokenRAYWHERE {
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

// parseInsert parses RAYERT RAYINTO table (columns) RAYVALUES (values)
func (p *Parser) parseInsert() (*Query, error) {
	query := &Query{Type: QueryInsert}

	// Expect RAYINTO
	if !p.expectPeek(TokenRAYINTO) {
		return nil, fmt.Errorf("expected RAYINTO")
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

	// Expect RAYVALUES
	if !p.expectPeek(TokenRAYVALUES) {
		return nil, fmt.Errorf("expected RAYVALUES")
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

// parseUpdate parses RAYDATE table RAYSET col=val [RAYWHERE condition]
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

	// Expect RAYSET
	if !p.expectPeek(TokenRAYSET) {
		return nil, fmt.Errorf("expected RAYSET")
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
	if p.peekToken.Type == TokenRAYWHERE {
		p.nextToken()
		cond, err := p.parseCondition()
		if err != nil {
			return nil, err
		}
		query.Conditions = cond
	}

	return query, nil
}

// parseDelete parses RAYLETE RAYFROM table [RAYWHERE condition]
func (p *Parser) parseDelete() (*Query, error) {
	query := &Query{Type: QueryDelete}

	// Expect RAYFROM
	if !p.expectPeek(TokenRAYFROM) {
		return nil, fmt.Errorf("expected RAYFROM")
	}

	// Expect table name
	if !p.expectPeek(TokenIdent) {
		return nil, fmt.Errorf("expected table name")
	}
	query.TableName = p.curToken.Literal

	// Optional WHERE
	if p.peekToken.Type == TokenRAYWHERE {
		p.nextToken()
		cond, err := p.parseCondition()
		if err != nil {
			return nil, err
		}
		query.Conditions = cond
	}

	return query, nil
}

// parseDrop parses RAYDROP RAYTABLE table
func (p *Parser) parseDrop() (*Query, error) {
	query := &Query{Type: QueryDrop}

	// Expect RAYTABLE
	if !p.expectPeek(TokenRAYTABLE) {
		return nil, fmt.Errorf("expected RAYTABLE")
	}

	// Expect table name
	if !p.expectPeek(TokenIdent) {
		return nil, fmt.Errorf("expected table name")
	}
	query.TableName = p.curToken.Literal

	return query, nil
}

// parseDescribe parses RAYC table
func (p *Parser) parseDescribe() (*Query, error) {
	query := &Query{Type: QueryDescribe}

	// Expect table name
	if !p.expectPeek(TokenIdent) {
		return nil, fmt.Errorf("expected table name")
	}
	query.TableName = p.curToken.Literal

	return query, nil
}

// parseShowTables parses RAYSHOW RAYTABLES
func (p *Parser) parseShowTables() (*Query, error) {
	query := &Query{Type: QueryShowTables}

	// Expect RAYTABLES
	if !p.expectPeek(TokenRAYTABLEs) {
		return nil, fmt.Errorf("expected RAYTABLES")
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
	if p.peekToken.Type == TokenRAYAND || p.peekToken.Type == TokenRAYOR {
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
	case TokenRAYLOVE:
		return "LIKE"
	case TokenRayamong:
		return "IN"
	case TokenRAYXIST:
		return "NOTNULL"
	case TokenRAYMISS:
		return "NULL"
	default:
		return p.curToken.Literal
	}
}
