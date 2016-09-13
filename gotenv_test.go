package gotenv

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var formats = []struct {
	in     string
	out    Env
	preset bool
}{
	// parses unquoted values
	{`FOO=bar`, Env{"FOO": "bar"}, false},

	// parses values with spaces around equal sign
	{`FOO =bar`, Env{"FOO": "bar"}, false},
	{`FOO= bar`, Env{"FOO": "bar"}, false},

	// parses double quoted values
	{`FOO="bar"`, Env{"FOO": "bar"}, false},

	// parses single quoted values
	{`FOO='bar'`, Env{"FOO": "bar"}, false},

	// parses escaped double quotes
	{`FOO="escaped\"bar"`, Env{"FOO": `escaped"bar`}, false},

	// parses empty values
	{`FOO=`, Env{"FOO": ""}, false},

	// expands variables found in values
	{"FOO=test\nBAR=$FOO", Env{"FOO": "test", "BAR": "test"}, false},

	// parses variables wrapped in brackets
	{"FOO=test\nBAR=${FOO}bar", Env{"FOO": "test", "BAR": "testbar"}, false},

	// reads variables from ENV when expanding if not found in local env
	{`BAR=$FOO`, Env{"BAR": "test"}, true},

	// expands undefined variables to an empty string
	{`BAR=$FOO`, Env{"BAR": ""}, false},

	// expands variables in quoted strings
	{"FOO=test\nBAR=\"quote $FOO\"", Env{"FOO": "test", "BAR": "quote test"}, false},

	// does not expand variables in single quoted strings
	{"BAR='quote $FOO'", Env{"BAR": "quote $FOO"}, false},

	// does not expand escaped variables
	{`FOO="foo\$BAR"`, Env{"FOO": "foo$BAR"}, false},
	{`FOO="foo\${BAR}"`, Env{"FOO": "foo${BAR}"}, false},
	{"FOO=test\nBAR=\"foo\\${FOO} ${FOO}\"", Env{"FOO": "test", "BAR": "foo${FOO} test"}, false},

	// parses yaml style options
	{"OPTION_A: 1", Env{"OPTION_A": "1"}, false},

	// parses export keyword
	{"export OPTION_A=2", Env{"OPTION_A": "2"}, false},

	// expands newlines in quoted strings
	{`FOO="bar\nbaz"`, Env{"FOO": "bar\nbaz"}, false},

	// parses variables with "." in the name
	{`FOO.BAR=foobar`, Env{"FOO.BAR": "foobar"}, false},

	// strips unquoted values
	{`foo=bar `, Env{"foo": "bar"}, false}, // not 'bar '

	// ignores empty lines
	{"\n \t  \nfoo=bar\n \nfizz=buzz", Env{"foo": "bar", "fizz": "buzz"}, false},

	// ignores inline comments
	{"foo=bar # this is foo", Env{"foo": "bar"}, false},

	// allows # in quoted value
	{`foo="bar#baz" # comment`, Env{"foo": "bar#baz"}, false},

	// ignores comment lines
	{"\n\n\n # HERE GOES FOO \nfoo=bar", Env{"foo": "bar"}, false},

	// parses # in quoted values
	{`foo="ba#r"`, Env{"foo": "ba#r"}, false},
	{"foo='ba#r'", Env{"foo": "ba#r"}, false},

	// supports carriage return
	{"FOO=bar\rbaz=fbb", Env{"FOO": "bar", "baz": "fbb"}, false},

	// supports carriage return combine with new line
	{"FOO=bar\r\nbaz=fbb", Env{"FOO": "bar", "baz": "fbb"}, false},

	// expands carriage return in quoted strings
	{`FOO="bar\rbaz"`, Env{"FOO": "bar\rbaz"}, false},
}

var fixtures = []struct {
	filename string
	results  Env
}{
	{
		"fixtures/exported.env",
		Env{
			"OPTION_A": "2",
			"OPTION_B": `\n`,
		},
	},
	{
		"fixtures/plain.env",
		Env{
			"OPTION_A": "1",
			"OPTION_B": "2",
			"OPTION_C": "3",
			"OPTION_D": "4",
			"OPTION_E": "5",
		},
	},
	{
		"fixtures/quoted.env",
		Env{
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
		Env{
			"OPTION_A": "1",
			"OPTION_B": "2",
			"OPTION_C": "",
			"OPTION_D": `\n`,
		},
	},
}

func TestParse(t *testing.T) {
	for _, tt := range formats {
		if tt.preset {
			os.Setenv("FOO", "test")
		}

		exp := Parse(strings.NewReader(tt.in))
		assert.Equal(t, tt.out, exp)
		os.Clearenv()
	}
}

func TestStrictParse(t *testing.T) {
	// incorrect line format
	env, err := StrictParse(strings.NewReader("lol$wut"))
	assert.NotNil(t, err)
	assert.Equal(t, Env{}, env)
}

func TestLoad(t *testing.T) {
	for _, tt := range fixtures {
		err := Load(tt.filename)
		assert.Nil(t, err)

		for key, val := range tt.results {
			assert.Equal(t, val, os.Getenv(key))
		}

		os.Clearenv()
	}
}

func TestLoad_default(t *testing.T) {
	k := "HELLO"
	v := "world"

	err := Load()
	assert.Nil(t, err)
	assert.Equal(t, v, os.Getenv(k))
	os.Clearenv()
}

func TestLoad_overriding(t *testing.T) {
	k := "HELLO"
	v := "universe"

	os.Setenv(k, v)
	err := Load()
	assert.Nil(t, err)
	assert.Equal(t, v, os.Getenv(k))
	os.Clearenv()
}

func TestLoad_invalidEnv(t *testing.T) {
	err := Load(".env.invalid")
	assert.NotNil(t, err)
}

func TestLoad_nonExist(t *testing.T) {
	file := ".env.not.exist"

	err := Load(file)
	if err == nil {
		t.Errorf("Load(`%s`) => error: `no such file or directory` != nil", file)
	}
}

func TestMustLoad(t *testing.T) {
	for _, tt := range fixtures {
		assert.NotPanics(t, func() {
			MustLoad(tt.filename)

			for key, val := range tt.results {
				assert.Equal(t, val, os.Getenv(key))
			}

			os.Clearenv()
		}, "Caling MustLoad should NOT panic")
	}
}

func TestMustLoad_default(t *testing.T) {
	assert.NotPanics(t, func() {
		MustLoad()

		tkey := "HELLO"
		val := "world"

		assert.Equal(t, val, os.Getenv(tkey))
		os.Clearenv()
	}, "Caling Load with no arguments should NOT panic")
}

func TestMustLoad_nonExist(t *testing.T) {
	assert.Panics(t, func() { MustLoad(".env.not.exist") }, "Caling MustLoad with non exist file SHOULD panic")
}

func TestOverLoad_overriding(t *testing.T) {
	k := "HELLO"
	v := "universe"

	os.Setenv(k, v)
	err := OverLoad()
	assert.Nil(t, err)
	assert.Equal(t, "world", os.Getenv(k))
	os.Clearenv()
}

func TestMustOverLoad_nonExist(t *testing.T) {
	assert.Panics(t, func() { MustOverLoad(".env.not.exist") }, "Caling MustOverLoad with non exist file SHOULD panic")
}

func TestApply(t *testing.T) {
	os.Setenv("HELLO", "world")
	r := strings.NewReader("HELLO=universe")
	err := Apply(r)
	assert.Nil(t, err)
	assert.Equal(t, "world", os.Getenv("HELLO"))
	os.Clearenv()
}

func TestOverApply(t *testing.T) {
	os.Setenv("HELLO", "world")
	r := strings.NewReader("HELLO=universe")
	err := OverApply(r)
	assert.Nil(t, err)
	assert.Equal(t, "universe", os.Getenv("HELLO"))
	os.Clearenv()
}
