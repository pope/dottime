package lexer

import (
	"pope/dottime/token"
)

type Lexer struct {
	input string
	pos   int
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) NextToken() token.Token {
	ch := l.skipIncidentalWhitespace()

	var tok token.Token
	switch ch {
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		return tok
	case '\n':
		l.charToken(&tok, token.NEWLINE)
		l.advance(1)
		return tok
	case '=':
		l.charToken(&tok, token.ASSIGN)
		l.advance(1)
		return tok
	case ':':
		l.charToken(&tok, token.COLON)
		l.advance(1)
		return tok
	case ',':
		l.charToken(&tok, token.COMMA)
		l.advance(1)
		return tok
	case ';':
		l.charToken(&tok, token.SEMICOLON)
		l.advance(1)
		return tok
	case '{':
		l.charToken(&tok, token.LBRACE)
		l.advance(1)
		return tok
	case '}':
		l.charToken(&tok, token.RBRACE)
		l.advance(1)
		return tok
	case '[':
		l.charToken(&tok, token.LBRACKET)
		l.advance(1)
		return tok
	case ']':
		l.charToken(&tok, token.RBRACKET)
		l.advance(1)
		return tok
	case '-':
		switch l.peek(1) {
		case '>':
			tok.Type = token.EDGEOP_DIRECTED
			tok.Literal = l.input[l.pos : l.pos+2]
			l.advance(2)
			return tok
		case '-':
			tok.Type = token.EDGEOP_UNIDIRECTED
			tok.Literal = l.input[l.pos : l.pos+2]
			l.advance(2)
			return tok
		default:
			if l.isNumeral() {
				l.consumeNumeral(&tok)
				return tok
			}
			goto ILLEGAL
		}
	case '#':
		switch l.prev() {
		case '\n', 0:
			l.consumeSingleLineComment(&tok)
			return tok
		default:
			goto ILLEGAL
		}
	case '"':
		l.consumeString(&tok)
		return tok
	case '/':
		switch l.peek(1) {
		case '/':
			l.consumeSingleLineComment(&tok)
			return tok
		default:
			// TODO(pope): Handle /* ... */ style comments
			goto ILLEGAL
		}
	default:
		if l.isNumeral() {
			l.consumeNumeral(&tok)
			return tok
		}
		if isAlpha(ch) {
			l.consumeIdent(&tok)
			return tok
		}
		// TODO(pope): Handle <...> HTML strings
		goto ILLEGAL
	}

ILLEGAL:
	l.charToken(&tok, token.ILLEGAL)
	return tok
}

func (l *Lexer) cur() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peek(n int) byte {
	pos := l.pos + n
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

func (l *Lexer) prev() byte {
	pos := l.pos - 1
	if pos < 0 {
		return 0
	}
	return l.input[pos]
}

func (l *Lexer) advance(n int) {
	l.pos += n
}

func (l *Lexer) skipIncidentalWhitespace() byte {
	for {
		ch := l.cur()
		switch ch {
		case ' ', '\t', '\r':
			l.advance(1)
		default:
			return ch
		}
	}
}

func (l *Lexer) consumeSingleLineComment(tok *token.Token) {
	pos := l.pos
loop:
	for {
		switch l.cur() {
		case '\n', 0:
			break loop
		default:
			l.advance(1)
		}
	}
	tok.Literal = l.input[pos:l.pos]
	tok.Type = token.COMMENT
}

func (l *Lexer) charToken(tok *token.Token, tt token.TokenType) {
	tok.Literal = l.input[l.pos : l.pos+1]
	tok.Type = tt
}

func (l *Lexer) consumeNumeral(tok *token.Token) {
	pos := l.pos

	cur := l.cur()
	if cur == '-' {
		l.advance(1)
		cur = l.cur()
	}

	for isDigit(cur) {
		l.advance(1)
		cur = l.cur()
	}

	if cur == '.' {
		l.advance(1)
		cur = l.cur()
	}

	for isDigit(cur) {
		l.advance(1)
		cur = l.cur()
	}

	tok.Literal = l.input[pos:l.pos]
	tok.Type = token.NUMERAL
}

func (l *Lexer) consumeIdent(tok *token.Token) {
	pos := l.pos

	cur := l.cur()
	for isAlpha(cur) || isDigit(cur) {
		l.advance(1)
		cur = l.cur()
	}

	tok.Literal = l.input[pos:l.pos]
	tok.Type = token.LookupIdentifier(tok.Literal)
}

func (l *Lexer) consumeString(tok *token.Token) {
	pos := l.pos

	// First character should be a " already.
	l.advance(1)
loop:
	for {
		switch l.cur() {
		case 0:
			tok.Literal = l.input[pos:l.pos]
			tok.Type = token.ILLEGAL
			return
		case '"':
			l.advance(1)
			break loop
		case '\\':
			l.advance(2)
		default:
			l.advance(1)
		}
	}

	tok.Literal = l.input[pos:l.pos]
	tok.Type = token.STRING
}

func (l *Lexer) isNumeral() bool {
	a := l.cur()
	if isDigit(a) {
		return true
	}
	if a == '-' {
		b := l.peek(1)
		if b == '.' {
			c := l.peek(2)
			return isDigit(c)
		}
		return isDigit(b)
	}
	return false
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		'\200' <= ch && ch <= '\377' ||
		ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
