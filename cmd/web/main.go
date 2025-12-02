package main

import (
	"flag"

	"booksonhooks.ca/internal/app"
	"booksonhooks.ca/internal/logger"
)

func main() {
	logger := logger.NewLogger("info")

	const port = ":4000"
	logger.Info("Starting application on port %s", port)

	addr := flag.String("addr", port, "HTTP network address")
	flag.Parse()

	srv := app.CreateApp(addr, &logger)
	err := srv.ListenAndServe()
	logger.ErrorLog.Fatal(err)
}
