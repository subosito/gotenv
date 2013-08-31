// package gotenv provides functionality to dynamically load the environment variables
package gotenv

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	linePattern     = `\A(?:export\s+)?([\w\.]+)(?:\s*=\s*|:\s+?)('(?:\'|[^'])*'|"(?:\"|[^"])*"|[^#\n]+)?(?:\s*\#.*)?\z`
	variablePattern = `(\\)?(\$)(\{?([A-Z0-9_]+)\}?)`
)

type Env map[string]string

// By default, it will load `.env` file on the current path and set the environment variables. You can supply filenames parameter to load your desired files.
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

		Parse(f)
	}

	return nil
}

func Parse(r io.Reader) Env {
	env := make(Env)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		et, ok := parseLine(scanner.Text())
		if ok {
			// set environment
			for key, val := range et {
				os.Setenv(key, val)
				env[key] = val
			}
		}
	}

	return env
}

func parseLine(s string) (Env, bool) {
	r := regexp.MustCompile(linePattern)
	matches := r.FindStringSubmatch(s)
	if len(matches) == 0 {
		return nil, false
	}

	key := matches[1]
	val := matches[2]

	// determine if string has quote prefix
	hq := strings.HasPrefix(val, `"`)

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

		if xv[1] == "\\" {
			replace = strings.Join(xv[2:4], "")
		} else {
			replace = os.Getenv(xv[4])
		}

		val = strings.Replace(val, strings.Join(xv[0:1], ""), replace, -1)
	}

	return Env{key: val}, true
}
