package parser_test

import (
	"bytes"
	"github.com/nais/naisplater/pkg/parser"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var parserTests = []struct {
	input  string
	output string
}{
	{
		`simple case`,
		`simple case`,
	},
	{
		`not s\n \no {{ simple .case }}`,
		`not s\n \no {{ simple .prefix.case }}`,
	},
	{
		`{ .single }`,
		`{ .single }`,
	},
	{
		`{ { .single }`,
		`{ { .single }`,
	},
	{
		`{{.single}}`,
		`{{.prefix.single}}`,
	},
	{
		`{{.multiple.levels.in.variable}}`,
		`{{.prefix.multiple.levels.in.variable}}`,
	},
	{
		`{{ range $i, $cluster := .clusters }}`,
		`{{ range $i, $cluster := .prefix.clusters }}`,
	},
	{
		`{{ (.paranthesis!=.foo) }}`,
		`{{ (.prefix.paranthesis!=.prefix.foo) }}`,
	},
}

func TestParser(t *testing.T) {

	for n, test := range parserTests {
		reader := strings.NewReader(test.input)
		writer := &bytes.Buffer{}
		prefix := ".prefix"

		t.Logf("### Test %d: '%s'", n+1, test.input)

		err := parser.ReplaceVariables(reader, writer, prefix)
		assert.NoError(t, err)
		assert.Equal(t, test.output, writer.String())
	}
}
