package scanner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Scanner struct {
	scan    *bufio.Scanner
	isAtEnd bool
	curLine []rune
	curChar rune
	linePos uint
	line    uint
	tokens  []Token
}

func New(fd io.Reader) Scanner {
	scan := bufio.NewScanner(fd)
	scan.Scan()

	firstLine := []rune(scan.Text())

	return Scanner{
		scan:    scan,
		curLine: firstLine,
		curChar: firstLine[0],
		line:    0,
		linePos: 0,
	}
}

func (s *Scanner) Lex() {
	for !s.isAtEnd {
		var tokenKind uint
		switch s.curChar {
		case '\n':
			s.addToken(NEWLINE, "")
			s.advanceLine()
			continue
		case '#':
			tokenKind = HASH
		default:
			var text strings.Builder

		out:
			for {
				switch s.curChar {
				case '\n', '#':
					break out
				}

				text.WriteRune(s.curChar)
				s.advance()
			}

			if text.Len() != 0 {
				s.addToken(TEXT, text.String())
			}

			continue
		}

		s.addToken(tokenKind, "")
		s.advance()
	}
}

func (s *Scanner) addToken(kind uint, value string) {
	pos := s.linePos
	if len(value) != 0 {
		pos = s.linePos - uint(len(value))
	}

	s.tokens = append(s.tokens, Token{
		Pos:   pos,
		Kind:  kind,
		Value: value,
		Line:  s.line,
	})
}

// advance the cursor
func (s *Scanner) advance() {
	if s.linePos+1 >= uint(len(s.curLine)) {
		s.curChar = '\n'
		s.linePos++
		return
	}

	s.linePos++
	s.curChar = s.curLine[s.linePos]
}

func (s *Scanner) advanceLine() {
	ok := s.scan.Scan()

	if !ok {
		if err := s.scan.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR reading input: ", s.scan.Err().Error())
		}

		s.isAtEnd = true
		return
	}

	s.curLine = []rune(s.scan.Text())
	s.line++
	s.linePos = 0

	// skip empty lines
	for len(s.curLine) == 0 {
		s.scan.Scan()
		s.curLine = []rune(s.scan.Text())
		s.line++
	}

	s.curChar = s.curLine[s.linePos]
}

func (s *Scanner) PrintTokens() {
	for _, token := range s.tokens {
		fmt.Printf("'%s', %d, %d, '%s'\n",
			TOKEN_LOOKUP[token.Kind],
			token.Pos,
			token.Line,
			token.Value,
		)
	}
}
