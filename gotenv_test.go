package gotenv

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var formats = []struct {
	in     string
	out    []map[string]string
	preset bool
}{
	// parses unquoted values
	{`FOO=bar`, []map[string]string{{"FOO": "bar"}}, false},

	// parses values with spaces around equal sign
	{`FOO =bar`, []map[string]string{{"FOO": "bar"}}, false},
	{`FOO= bar`, []map[string]string{{"FOO": "bar"}}, false},

	// parses double quoted values
	{`FOO="bar"`, []map[string]string{{"FOO": "bar"}}, false},

	// parses single quoted values
	{`FOO='bar'`, []map[string]string{{"FOO": "bar"}}, false},

	// parses escaped double quotes
	{`FOO="escaped\"bar"`, []map[string]string{{"FOO": `escaped"bar`}}, false},

	// parses empty values
	{`FOO=`, []map[string]string{{"FOO": ""}}, false},

	// expands variables found in values
	{"FOO=test\nBAR=$FOO", []map[string]string{{"FOO": "test"}, {"BAR": "test"}}, false},

	// parses variables wrapped in brackets
	{"FOO=test\nBAR=${FOO}bar", []map[string]string{{"FOO": "test"}, {"BAR": "testbar"}}, false},

	// reads variables from ENV when expanding if not found in local env
	{`BAR=$FOO`, []map[string]string{{"BAR": "test"}}, true},

	// expands undefined variables to an empty string
	{`BAR=$FOO`, []map[string]string{{"BAR": ""}}, false},

	// expands variables in quoted strings
	{"FOO=test\nBAR='quote $FOO'", []map[string]string{{"FOO": "test"}, {"BAR": "quote test"}}, false},

	// does not expand escaped variables
	{`FOO="foo\$BAR"`, []map[string]string{{"FOO": "foo$BAR"}}, false},
	{`FOO="foo\${BAR}"`, []map[string]string{{"FOO": "foo${BAR}"}}, false},

	// parses yaml style options
	{"OPTION_A: 1", []map[string]string{{"OPTION_A": "1"}}, false},

	// parses export keyword
	{"export OPTION_A=2", []map[string]string{{"OPTION_A": "2"}}, false},

	// expands newlines in quoted strings
	{`FOO="bar\nbaz"`, []map[string]string{{"FOO": "bar\nbaz"}}, false},

	// parses varibales with "." in the name
	{`FOO.BAR=foobar`, []map[string]string{{"FOO.BAR": "foobar"}}, false},

	// strips unquoted values
	{`foo=bar `, []map[string]string{{"foo": "bar"}}, false}, // not 'bar '

	// ignores empty lines
	{"\n \t  \nfoo=bar\n \nfizz=buzz", []map[string]string{{"foo": "bar"}, {"fizz": "buzz"}}, false},

	// ignores inline comments
	{"foo=bar # this is foo", []map[string]string{{"foo": "bar"}}, false},

	// allows # in quoted value
	{`foo="bar#baz" # comment`, []map[string]string{{"foo": "bar#baz"}}, false},

	// ignores comment lines
	{"\n\n\n # HERE GOES FOO \nfoo=bar", []map[string]string{{"foo": "bar"}}, false},

	// parses # in quoted values
	{`foo="ba#r"`, []map[string]string{{"foo": "ba#r"}}, false},
	{"foo='ba#r'", []map[string]string{{"foo": "ba#r"}}, false},

	// incorrect line format
	{"lol$wut", []map[string]string{}, false},
}

var fixtures = []struct {
	filename string
	results  map[string]string
}{
	{
		"fixtures/exported.env",
		map[string]string{
			"OPTION_A": "2",
			"OPTION_B": `\n`,
		},
	},
	{
		"fixtures/plain.env",
		map[string]string{
			"OPTION_A": "1",
			"OPTION_B": "2",
			"OPTION_C": "3",
			"OPTION_D": "4",
			"OPTION_E": "5",
		},
	},
	{
		"fixtures/quoted.env",
		map[string]string{
			"OPTION_A": "1",
			"OPTION_B": "2",
			"OPTION_C": "",
			"OPTION_D": `\n`,
			"OPTION_E": "1",
			"OPTION_F": "2",
			"OPTION_G": "",
			"OPTION_H": "\n",
		},
	},
	{
		"fixtures/yaml.env",
		map[string]string{
			"OPTION_A": "1",
			"OPTION_B": "2",
			"OPTION_C": "",
			"OPTION_D": `\n`,
		},
	},
}

func TestParse(t *testing.T) {
	for i, tt := range formats {

		// reset environments
		os.Clearenv()

		if tt.preset {
			os.Setenv("FOO", "test")
		}

		exp := Parse(strings.NewReader(tt.in))

		x := fmt.Sprintf("%+v\n", exp)
		o := fmt.Sprintf("%+v\n", tt.out)

		if x != o {
			t.Logf("%q\n", tt.in)
			t.Errorf("(%d) %s != %s\n", i, x, o)
		}
	}
}

func TestLoad(t *testing.T) {
	for i, tt := range fixtures {
		Load(tt.filename)

		for key, val := range tt.results {
			if eval := os.Getenv(key); eval != val {
				t.Errorf("(%d) %s => %s != %s", i, key, eval, val)
			}
		}

		// reset environments
		os.Clearenv()
	}
}
