package lexer_test

import (
	"github.com/nais/naisplater/pkg/lexer"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type result struct {
	class int
	str   string
}

var lexerTests = []struct {
	input    string
	expected []result
}{
	{
		`a    normal sentence`,
		[]result{
			{class: lexer.TokenIdentifier, str: `a`},
			{class: lexer.TokenWhitespace, str: `    `},
			{class: lexer.TokenIdentifier, str: `normal`},
			{class: lexer.TokenWhitespace, str: ` `},
			{class: lexer.TokenIdentifier, str: `sentence`},
			{class: lexer.TokenEnd, str: ``},
		},
	},
	{
		`{{ this is .variable }}`,
		[]result{
			{class: lexer.TokenCurlyLeft, str: `{{`},
			{class: lexer.TokenWhitespace, str: ` `},
			{class: lexer.TokenIdentifier, str: `this`},
			{class: lexer.TokenWhitespace, str: ` `},
			{class: lexer.TokenIdentifier, str: `is`},
			{class: lexer.TokenWhitespace, str: ` `},
			{class: lexer.TokenIdentifier, str: `.variable`},
			{class: lexer.TokenWhitespace, str: ` `},
			{class: lexer.TokenCurlyRight, str: `}}`},
			{class: lexer.TokenEnd, str: ``},
		},
	},
}

func TestLexer(t *testing.T) {

	for n, test := range lexerTests {

		index := 0
		reader := strings.NewReader(test.input)
		scanner := lexer.NewScanner(reader)

		t.Logf("### Test %d: '%s'", n+1, test.input)

		for {
			class, str := scanner.Scan()

			if index == len(test.expected) {
				if class == lexer.TokenEnd {
					break
				}
				t.Fatalf("Tokenizer generated too many tokens!")
			}

			t.Logf("Token %d: class='%d', literal='%s'", index, class, str)

			check := test.expected[index]

			assert.Equal(t, check.class, class,
				"Token class for token %d is wrong; expected %d but got %d", index, check.class, class)
			assert.Equal(t, check.str, str,
				"String check against token %d failed; expected '%s' but got '%s'", index, check.str, str)

			index++
		}
	}
}
