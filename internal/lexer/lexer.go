// Package lexer provides tokenization for GitHub Actions DSL
package lexer

import (
	"fmt"
)

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Identifiers and literals
	IDENT  // workflow, job, step, variable names
	STRING // "string literals"
	NUMBER // numeric literals

	// Keywords
	WORKFLOW
	JOB
	STEP
	ON
	USES
	RUN
	WITH
	ENV
	NEEDS
	IF
	STRATEGY
	MATRIX

	// Loop keywords
	FOR
	FOREACH
	REPEAT
	IN
	RANGE

	// Operators and delimiters
	ASSIGN    // =
	LBRACE    // {
	RBRACE    // }
	LPAREN    // (
	RPAREN    // )
	LBRACKET  // [
	RBRACKET  // ]
	COMMA     // ,
	SEMICOLON // ;
	NEWLINE   // \n

	// Additional operators for future expansion
	DOT    // .
	COLON  // :
	PIPE   // |
	DOTDOT // .. (for ranges)
)

// Token represents a lexical token
type Token struct {
	Type     TokenType
	Literal  string
	Line     int
	Column   int
	Position int
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %s, Literal: %s, Line: %d, Col: %d}",
		t.Type.String(), t.Literal, t.Line, t.Column)
}

// TokenType.String returns string representation of token type
func (tt TokenType) String() string {
	switch tt {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case WORKFLOW:
		return "WORKFLOW"
	case JOB:
		return "JOB"
	case STEP:
		return "STEP"
	case ON:
		return "ON"
	case USES:
		return "USES"
	case RUN:
		return "RUN"
	case WITH:
		return "WITH"
	case ENV:
		return "ENV"
	case NEEDS:
		return "NEEDS"
	case IF:
		return "IF"
	case STRATEGY:
		return "STRATEGY"
	case MATRIX:
		return "MATRIX"
	case FOR:
		return "FOR"
	case FOREACH:
		return "FOREACH"
	case REPEAT:
		return "REPEAT"
	case IN:
		return "IN"
	case RANGE:
		return "RANGE"
	case ASSIGN:
		return "ASSIGN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case COMMA:
		return "COMMA"
	case SEMICOLON:
		return "SEMICOLON"
	case NEWLINE:
		return "NEWLINE"
	case DOT:
		return "DOT"
	case COLON:
		return "COLON"
	case PIPE:
		return "PIPE"
	case DOTDOT:
		return "DOTDOT"
	default:
		return "UNKNOWN"
	}
}

// Keywords maps keyword strings to their token types
var keywords = map[string]TokenType{
	"workflow": WORKFLOW,
	"job":      JOB,
	"step":     STEP,
	"on":       ON,
	"uses":     USES,
	"run":      RUN,
	"with":     WITH,
	"env":      ENV,
	"needs":    NEEDS,
	"if":       IF,
	"strategy": STRATEGY,
	"matrix":   MATRIX,
	"for":      FOR,
	"foreach":  FOREACH,
	"repeat":   REPEAT,
	"in":       IN,
	"range":    RANGE,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// Lexer represents the lexical analyzer
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number
	column       int  // current column number
}

// New creates a new lexer instance
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL character represents EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace skips whitespace characters except newlines
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// skipComment skips single-line comments starting with //
func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '-' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readString reads a quoted string literal
func (l *Lexer) readString() string {
	position := l.position + 1 // skip opening quote
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
		// Handle escaped quotes
		if l.ch == '\\' && l.peekChar() == '"' {
			l.readChar() // skip backslash
			l.readChar() // skip escaped quote
		}
	}
	return l.input[position:l.position]
}

// NextToken scans and returns the next token
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	// Handle comments
	if l.ch == '/' && l.peekChar() == '/' {
		l.skipComment()
		return l.NextToken()
	}

	tok.Line = l.line
	tok.Column = l.column
	tok.Position = l.position

	switch l.ch {
	case '=':
		tok = Token{Type: ASSIGN, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '{':
		tok = Token{Type: LBRACE, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '[':
		tok = Token{Type: LBRACKET, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ']':
		tok = Token{Type: RBRACKET, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ',':
		tok = Token{Type: COMMA, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ';':
		tok = Token{Type: SEMICOLON, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '\n':
		tok = Token{Type: NEWLINE, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '.':
		if l.peekChar() == '.' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: DOTDOT, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = Token{Type: DOT, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	case ':':
		tok = Token{Type: COLON, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '|':
		tok = Token{Type: PIPE, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok // early return to avoid readChar() call
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	}

	l.readChar()
	return tok
}

// readNumber reads a number (integers for simplicity)
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// isLetter checks if character is a letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit checks if character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// GetAllTokens returns all tokens from the input (useful for debugging)
func (l *Lexer) GetAllTokens() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return tokens
}

// Error represents a lexical error
type Error struct {
	Position int
	Line     int
	Column   int
	Message  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("lexer error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}
