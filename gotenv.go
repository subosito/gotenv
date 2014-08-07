// Package gotenv provides functionality to dynamically load the environment variables
package gotenv

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	// Pattern for detecting valid line format
	linePattern = `\A(?:export\s+)?([\w\.]+)(?:\s*=\s*|:\s+?)('(?:\'|[^'])*'|"(?:\"|[^"])*"|[^#\n]+)?(?:\s*\#.*)?\z`

	// Pattern for detecting valid variable within a value
	variablePattern = `(\\)?(\$)(\{?([A-Z0-9_]+)\}?)`
)

// Holds key/value pair of valid environment variable
type Env map[string]string

/*
Load is function to load a file or multiple files and then export the valid variables which found into environment variables.
When it's called with no argument, it will load `.env` file on the current path and set the environment variables.
Otherwise, it will loop over the filenames parameter and set the proper environment variables.

	// processing `.env`
	gotenv.Load()

	// processing multiple files
	gotenv.Load("production.env", "credentials")

*/
func Load(filenames ...string) error {
	if len(filenames) == 0 {
		filenames = []string{".env"}
	}

	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		// set environment
		env := Parse(f)
		for key, val := range env {
			os.Setenv(key, val)
		}
	}

	return nil
}

// Parse is a function to parse line by line any io.Reader supplied and returns the valid Env key/value pair of valid variables.
// It expands the value of a variable from environment variable, but does not set the value to the environment itself.
// This function is skipping any invalid lines and only processing the valid one.
func Parse(r io.Reader) Env {
	env, _ := StrictParse(r)
	return env
}

// StrictParse is a function to parse line by line any io.Reader supplied and returns the valid Env key/value pair of valid variables.
// It expands the value of a variable from environment variable, but does not set the value to the environment itself.
// This function is returning an error if there is any invalid lines.
func StrictParse(r io.Reader) (Env, error) {
	env := make(Env)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		err := parseLine(scanner.Text(), env)
		if err != nil {
			return env, err
		}
	}

	return env, nil
}

func parseLine(s string, env Env) error {
	r := regexp.MustCompile(linePattern)
	matches := r.FindStringSubmatch(s)
	if len(matches) == 0 {
		st := strings.TrimSpace(s)

		if (st == "") || strings.HasPrefix(st, "#") {
			return nil
		}

		return errors.New(fmt.Sprintf("Line `%s` doesn't match format", s))
	}

	key := matches[1]
	val := matches[2]

	// determine if string has quote prefix
	hq := strings.HasPrefix(val, `"`)

	// determine if string has single quote prefix
	hs := strings.HasPrefix(val, `'`)

	// trim whitespace
	val = strings.Trim(val, " ")

	// remove quotes '' or ""
	rq := regexp.MustCompile(`\A(['"])(.*)(['"])\z`)
	val = rq.ReplaceAllString(val, "$2")

	if hq {
		val = strings.Replace(val, `\n`, "\n", -1)
		// Unescape all characters except $ so variables can be escaped properly
		re := regexp.MustCompile(`\\([^$])`)
		val = re.ReplaceAllString(val, "$1")
	}

	rv := regexp.MustCompile(variablePattern)
	xv := rv.FindStringSubmatch(val)

	if len(xv) > 0 {
		var replace string
		var ok bool

		if xv[1] == "\\" {
			replace = strings.Join(xv[2:4], "")
		} else {
			replace, ok = env[xv[4]]
			if !ok {
				replace = os.Getenv(xv[4])
			}
		}

		if !hs {
			val = strings.Replace(val, strings.Join(xv[0:1], ""), replace, -1)
		}
	}

	env[key] = val
	return nil
}
