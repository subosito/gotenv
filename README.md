# Gotenv

Loads environment variables from `.env` file.

## Usage

You can go get this package as usual:

```bash
$ go get github.com/subosito/gotenv
```

Add your application configuration to `.env` file on your root directory of your project:

```
APP_ID=1234567
API_SECRET=abcdef
```

gotenv support various format for the `.env` file, you can see more formats on fixtures folder.

Then on your application code, put:

```go
gotenv.Load()
```

It will then load `.env` file and export the valid variables to the environment variables.

You can also load other than `.env` file if you wish.

```go
gotenv.Load("production.env", "credentials")
```

That's it :)

## TODO

- Write proper documentation

## Notes

Since `gotenv` is a Go port of [`dotenv`](https://github.com/bkeepers/dotenv) project, most logic and regexp pattern is taken from there and it will be compatible as close as possible with the `dotenv`.

