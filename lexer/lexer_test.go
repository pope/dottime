package lexer

import (
	"testing"

	"pope/dottime/token"
)

type tokenTest struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func verifyToken(t *testing.T, name string, i int, tt *tokenTest, tok *token.Token) {
	if tok.Type != tt.expectedType {
		t.Fatalf("%s[%d] - token type wrong. expected=%q, got=%q",
			name, i, tt.expectedType, tok.Type)
	}
	if tok.Literal != tt.expectedLiteral {
		t.Fatalf("%s[%d] - literal wrong. expected=%q, got=%q",
			name, i, tt.expectedLiteral, tok.Literal)
	}
}

func TestSimpleDirectional(t *testing.T) {
	input := `digraph graphname {
  a -> b -> c;
	b -> d;
}
`
	tests := []tokenTest{
		{token.DIGRAPH, "digraph"},
		{token.IDENT, "graphname"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},

		{token.IDENT, "a"},
		{token.EDGEOP_DIRECTED, "->"},
		{token.IDENT, "b"},
		{token.EDGEOP_DIRECTED, "->"},
		{token.IDENT, "c"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},

		{token.IDENT, "b"},
		{token.EDGEOP_DIRECTED, "->"},
		{token.IDENT, "d"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},

		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		verifyToken(t, "TestSimpleDirectional", i, &tt, &tok)
	}
}

func TestSimpleUnidirectional(t *testing.T) {
	input := `// The graph name and semicolons are optional
graph graphname {
  a -- b -- c;
	b -- d;
}
`
	tests := []tokenTest{
		{token.COMMENT, "// The graph name and semicolons are optional"},
		{token.NEWLINE, "\n"},
		{token.GRAPH, "graph"},
		{token.IDENT, "graphname"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},

		{token.IDENT, "a"},
		{token.EDGEOP_UNIDIRECTED, "--"},
		{token.IDENT, "b"},
		{token.EDGEOP_UNIDIRECTED, "--"},
		{token.IDENT, "c"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},

		{token.IDENT, "b"},
		{token.EDGEOP_UNIDIRECTED, "--"},
		{token.IDENT, "d"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},

		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		verifyToken(t, "TestSimpleUnidirectional", i, &tt, &tok)
	}
}

func TestAttributes(t *testing.T) {
	input := `graph graphname {
    // This attribute applies to the graph itself
    size="1,1";
    // The label attribute can be used to change the label of a node
    a [label="Foo"];
    // Here, the node shape is changed.
    b [shape=box];
    // These edges both have different line properties
    a -- b -- c [color=blue];
    b -- d [style=dotted];
    // [style=invis] hides a node.
}
`
	tests := []tokenTest{
		{token.GRAPH, "graph"},
		{token.IDENT, "graphname"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},

		{token.COMMENT, "// This attribute applies to the graph itself"},
		{token.NEWLINE, "\n"},

		{token.IDENT, "size"},
		{token.ASSIGN, "="},
		{token.STRING, "\"1,1\""},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},

		{token.COMMENT, "// The label attribute can be used to change the label of a node"},
		{token.NEWLINE, "\n"},

		{token.IDENT, "a"},
		{token.LBRACKET, "["},
		{token.IDENT, "label"},
		{token.ASSIGN, "="},
		{token.STRING, "\"Foo\""},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},

		{token.COMMENT, "// Here, the node shape is changed."},
		{token.NEWLINE, "\n"},

		{token.IDENT, "b"},
		{token.LBRACKET, "["},
		{token.IDENT, "shape"},
		{token.ASSIGN, "="},
		{token.IDENT, "box"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},

		{token.COMMENT, "// These edges both have different line properties"},
		{token.NEWLINE, "\n"},

		{token.IDENT, "a"},
		{token.EDGEOP_UNIDIRECTED, "--"},
		{token.IDENT, "b"},
		{token.EDGEOP_UNIDIRECTED, "--"},
		{token.IDENT, "c"},
		{token.LBRACKET, "["},
		{token.IDENT, "color"},
		{token.ASSIGN, "="},
		{token.IDENT, "blue"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},

		{token.IDENT, "b"},
		{token.EDGEOP_UNIDIRECTED, "--"},
		{token.IDENT, "d"},
		{token.LBRACKET, "["},
		{token.IDENT, "style"},
		{token.ASSIGN, "="},
		{token.IDENT, "dotted"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},

		{token.COMMENT, "// [style=invis] hides a node."},
		{token.NEWLINE, "\n"},

		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		verifyToken(t, "TestAttributes", i, &tt, &tok)
	}
}

func TestOctothorpeComment(t *testing.T) {
	input := `# Hello world
# What is up
	# Illegal because not at start of newline`

	tests := []tokenTest{
		{token.COMMENT, "# Hello world"},
		{token.NEWLINE, "\n"},
		{token.COMMENT, "# What is up"},
		{token.NEWLINE, "\n"},
		{token.ILLEGAL, "#"},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		verifyToken(t, "TestOctothorpeComment", i, &tt, &tok)
	}
}
