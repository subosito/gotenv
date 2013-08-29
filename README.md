# gotenv

Dynamic way to load environment variables in Go.

|-              | -                                                  |
|---------------|----------------------------------------------------|
| Build Status  | [![Build Status][travis-img]][travis-url]          |
| Coverage      | [![Coverage Status][coveralls-img]][coveralls-url] |
| Documentation | http://godoc.org/github.com/subosito/gotenv        |

## Usage

```go
import "github.com/subosito/gotenv"
```

The gotenv supports two ways for loading environment variables:

1. Loading from a (.env) file
2. Parsing from a string.

### Loading from a (.env) file

Add your configuration to `.env` file on your root directory of your project:

```
APP_ID=1234567
API_SECRET=abcdef
```

Then on your application code, put:

```go
gotenv.Load()
```

Behind the scene it will then load `.env` file and export the valid variables to the environment variables. Make sure you call the method as soon as possible to ensure all variables are loaded.

You can also load other than `.env` file if you wish. Just supply filenames when calling `Load()`:

```go
gotenv.Load("production.env", "credentials")
```

That's it :)

### Parsing from a string

Besides loading from file, gotenv also support parsing environment variables from a string. gotenv provides functions for that purpose:

- for single line string:

```go
pair, ok := gotenv.ParseLine(`FOO="bar"`)
if ok {
	fmt.Printf("%+v\n", pair) // map[string]string{"FOO": "bar"}
}
```

- for multiline string _(can use for single line string too)_:

```go
// import "strings"

pairs = gotenv.Parse(strings.NewReader("FOO=test\nBAR=$FOO"))
fmt.Printf("%+v\n", pairs) // []map[string]string{{"FOO": "test"}, {"BAR": "test"}}

pairs := gotenv.Parse(strings.NewReader(`FOO="bar"`))
fmt.Printf("%+v\n", pairs) // []map[string]string{{"FOO": "bar"}}
```

## Formats

The gotenv supports various format for the `.env` file, you can see more formats on [fixtures](./fixtures) folder.

## TODO

- Write proper documentation

## Notes

Since `gotenv` is a Go port of [`dotenv`](https://github.com/bkeepers/dotenv) project, most logic and regexp pattern is taken from there and it will be compatible as close as possible with the `dotenv`.

[travis-img]: https://travis-ci.org/subosito/gotenv.png
[travis-url]: https://travis-ci.org/subosito/gotenv
[coveralls-img]: https://coveralls.io/repos/subosito/gotenv/badge.png?branch=master
[coveralls-url]: https://coveralls.io/r/subosito/gotenv?branch=master

