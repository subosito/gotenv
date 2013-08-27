package gotenv_test

import (
	"github.com/subosito/gotenv"
)

func ExampleLoad() {
	// Load default .env file
	gotenv.Load()

	// Load particular files
	gotenv.Load("production.env", "credentials")
}
