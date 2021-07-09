package parser

import (
	"fmt"
	"github.com/nais/naisplater/pkg/lexer"
	"io"
	"strings"
)

func ReplaceVariables(r io.Reader, w io.Writer, prefix string) error {
	var insideExpression bool

	scan := lexer.NewScanner(r)

	for {
		tok, lit := scan.Scan()
		if tok == lexer.TokenEnd {
			return nil
		}

		if insideExpression && strings.HasPrefix(lit, ".") {
			lit = prefix + lit
		}

		_, err := w.Write([]byte(lit))
		if err != nil {
			return err
		}

		if tok == lexer.TokenCurlyLeft && len(lit) == 2 {
			if insideExpression {
				return fmt.Errorf("template error: double nested expression")
			}
			insideExpression = true
		}

		if tok == lexer.TokenCurlyRight && len(lit) == 2 {
			if !insideExpression {
				return fmt.Errorf("template error: end expression, but not inside expression")
			}
			insideExpression = false
		}
	}
}
