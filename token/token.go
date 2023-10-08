package token

import (
	"strings"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	STRICT   TokenType = "STRICT"
	GRAPH    TokenType = "GRAPH"
	DIGRAPH  TokenType = "DIGRAPH"
	NODE     TokenType = "NODE"
	EDGE     TokenType = "EDGE"
	SUBGRAPH TokenType = "SUBGRAPH"

	IDENT   TokenType = "IDENT"
	STRING  TokenType = "STRING"
	NUMERAL TokenType = "NUMERAL"

	ASSIGN    TokenType = "="
	COLON     TokenType = ":"
	SEMICOLON TokenType = ";"
	COMMA     TokenType = ","
	NEWLINE   TokenType = "NEWLINE"

	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	EDGEOP_DIRECTED    TokenType = "->"
	EDGEOP_UNIDIRECTED TokenType = "--"

	COMMENT TokenType = "COMMENT"
)

func LookupIdentifier(ident string) TokenType {
	switch strings.ToLower(ident) {
	case "strict":
		return STRICT
	case "graph":
		return GRAPH
	case "digraph":
		return DIGRAPH
	case "node":
		return NODE
	case "edge":
		return EDGE
	case "subgraph":
		return SUBGRAPH
	default:
		return IDENT
	}
}
