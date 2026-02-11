package parser

import (
	"strings"
	"unicode"
)

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	TokenEOF TokenType = iota
	TokenIllegal

	// Identifiers and literals
	TokenIdent
	TokenString
	TokenNumber
	TokenBoolean

	// Friska Keywords
	TokenFrisrate // FRISRATE (CREATE)
	TokenFrislect // FRISLECT (SELECT)
	TokenFrisert  // FRISERT (INSERT)
	TokenFrisdate // FRISDATE (UPDATE)
	TokenFrislete // FRISLETE (DELETE)
	TokenFrisdrop // FRISDROP (DROP)
	TokenFrisc    // FRISC (DESCRIBE)
	TokenFrisshow // FRISSHOW (SHOW)

	// Friska Clauses
	TokenFriskable  // FRISKABLE (TABLE)
	TokenFriskables // FRISKABLES (TABLES)
	TokenFrisfrom   // FRISFROM (FROM)
	TokenFriswhere  // FRISWHERE (WHERE)
	TokenFrisset    // FRISSET (SET)
	TokenFrisinto   // FRISINTO (INTO)
	TokenFrisvalues // FRISVALUES (VALUES)

	// Logical Operators
	TokenFrisand // FRISAND (AND)
	TokenFrisor  // FRISOR (OR)
	TokenFrisnot // FRISNOT (NOT)

	// Special Operators
	TokenFrislove  // FRISLOVE (LIKE)
	TokenFrisamong // FRISAMONG (IN)
	TokenFrisxist  // FRISXIST (IS NOT NULL)
	TokenFrismiss  // FRISMISS (IS NULL)

	// Comparison Operators
	TokenAbove   // ABOVE (>)
	TokenBelow   // BELOW (<)
	TokenAtleast // ATLEAST (>=)
	TokenAtmost  // ATMOST (<=)

	// Symbols
	TokenComma
	TokenLeftParen
	TokenRightParen
	TokenAsterisk
	TokenEquals
	TokenNotEquals
	TokenGreater
	TokenLess
	TokenGreaterEq
	TokenLessEq
)

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Lexer tokenizes Friska queries
type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
	line    int
	column  int
}

// NewLexer creates a new lexer instance
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// NextToken returns the next token
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case 0:
		tok.Type = TokenEOF
		tok.Literal = ""
	case ',':
		tok.Type = TokenComma
		tok.Literal = string(l.ch)
	case '(':
		tok.Type = TokenLeftParen
		tok.Literal = string(l.ch)
	case ')':
		tok.Type = TokenRightParen
		tok.Literal = string(l.ch)
	case '*':
		tok.Type = TokenAsterisk
		tok.Literal = string(l.ch)
	case '=':
		tok.Type = TokenEquals
		tok.Literal = string(l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenNotEquals
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenGreaterEq
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = TokenGreater
			tok.Literal = string(l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = TokenLessEq
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = TokenLess
			tok.Literal = string(l.ch)
		}
	case '\'', '"':
		tok.Type = TokenString
		tok.Literal = l.readString(l.ch)
	case '%':
		// Part of LIKE pattern, treat as string
		tok.Type = TokenString
		tok.Literal = l.readPattern()
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = TokenNumber
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok.Type = TokenIllegal
			tok.Literal = string(l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	l.column++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readString(quote byte) string {
	pos := l.pos + 1
	for {
		l.readChar()
		if l.ch == quote || l.ch == 0 {
			break
		}
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readPattern() string {
	pos := l.pos
	for l.ch == '%' || isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

// lookupIdent checks if identifier is a keyword
func lookupIdent(ident string) TokenType {
	keywords := map[string]TokenType{
		// Core commands
		"FRISRATE": TokenFrisrate,
		"FRISLECT": TokenFrislect,
		"FRISERT":  TokenFrisert,
		"FRISDATE": TokenFrisdate,
		"FRISLETE": TokenFrislete,
		"FRISDROP": TokenFrisdrop,
		"FRISC":    TokenFrisc,
		"FRISSHOW": TokenFrisshow,

		// Clauses
		"FRISKABLE":  TokenFriskable,
		"FRISKABLES": TokenFriskables,
		"FRISFROM":   TokenFrisfrom,
		"FRISWHERE":  TokenFriswhere,
		"FRISSET":    TokenFrisset,
		"FRISINTO":   TokenFrisinto,
		"FRISVALUES": TokenFrisvalues,

		// Logical
		"FRISAND": TokenFrisand,
		"FRISOR":  TokenFrisor,
		"FRISNOT": TokenFrisnot,

		// Special operators
		"FRISLOVE":  TokenFrislove,
		"FRISAMONG": TokenFrisamong,
		"FRISXIST":  TokenFrisxist,
		"FRISMISS":  TokenFrismiss,

		// Comparison
		"ABOVE":   TokenAbove,
		"BELOW":   TokenBelow,
		"ATLEAST": TokenAtleast,
		"ATMOST":  TokenAtmost,

		// Booleans
		"true":  TokenBoolean,
		"false": TokenBoolean,
		"TRUE":  TokenBoolean,
		"FALSE": TokenBoolean,
	}

	upper := strings.ToUpper(ident)
	if tok, ok := keywords[upper]; ok {
		return tok
	}
	return TokenIdent
}
