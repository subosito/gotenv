# gotenv

[![Build Status](https://travis-ci.org/subosito/gotenv.svg?branch=master)](https://travis-ci.org/subosito/gotenv)
[![Build status](https://ci.appveyor.com/api/projects/status/wb2e075xkfl0m0v2/branch/master?svg=true)](https://ci.appveyor.com/project/subosito/gotenv/branch/master)
[![Coverage Status](https://badgen.net/codecov/c/github/subosito/gotenv)](https://codecov.io/gh/subosito/gotenv)
[![Go Report Card](https://goreportcard.com/badge/github.com/subosito/gotenv)](https://goreportcard.com/report/github.com/subosito/gotenv)
[![GoDoc](https://godoc.org/github.com/subosito/gotenv?status.svg)](https://godoc.org/github.com/subosito/gotenv)

Load environment variables dynamically in Go.

## Installation

```bash
$ go get github.com/subosito/gotenv
```

## Usage

Put the gotenv package on your `import` statement:

```go
import "github.com/subosito/gotenv"
```

By default, `Load` will look for a file called `.env` in the current working directory.

```go
gotenv.Load()
```

Behind the scene it will then load `.env` file and export the valid variables to the environment variables. Make sure you call the method as soon as possible to ensure it loads all variables, say, put it on `init()` function.

Once loaded you can use `os.Getenv()` to get the value of the variable.

Let's say you have `.env` file:

```
APP_ID=1234567
APP_SECRET=abcdef
```

Here's the example of your app:

```go
package main

import (
	"github.com/subosito/gotenv"
	"log"
	"os"
)

func init() {
	gotenv.Load()
}

func main() {
	log.Println(os.Getenv("APP_ID"))     // "1234567"
	log.Println(os.Getenv("APP_SECRET")) // "abcdef"
}
```

You can also load other than `.env` file if you wish. Just supply filenames when calling `Load()`. It will load them in order and the first value set for a variable will win.:

```go
gotenv.Load(".env.production", "credentials")
```

That's it :)

### Another Scenario

Just in case you want to parse environment variables from any `io.Reader`, gotenv keeps its `Parse()` function as public API so you can use that.

```go
// import "strings"

pairs := gotenv.Parse(strings.NewReader("FOO=test\nBAR=$FOO"))
// gotenv.Env{"FOO": "test", "BAR": "test"}

pairs = gotenv.Parse(strings.NewReader(`FOO="bar"`))
// gotenv.Env{"FOO": "bar"}
```

Parse ignores invalid lines and returns `Env` of valid environment variables.

## Notes

The gotenv package is a Go port of [`dotenv`](https://github.com/bkeepers/dotenv) project. It aims to be compatible as close as possible.
