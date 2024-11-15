package main

import (
	"fmt"
	"strings"

	"github.com/antony-with-no-h/ds-nq-markdown-parser/scanner"
)

func main() {
	input := strings.NewReader(`# Heading 1
paragraph`)

	// scanner.New()
	s := scanner.New(input)

	// scanner.Lex()
	s.Lex()

	fmt.Printf("%+v\n", s)

}
