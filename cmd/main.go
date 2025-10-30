package main

import (
	"log/slog"

	"github.com/simonnix/fiber-csrf-v3-test"
)

func main() {
	app := csrfv3.GetApp()
	if err := app.Listen(":3000"); err != nil {
		slog.Error("Server returned: " + err.Error())
	}
}
