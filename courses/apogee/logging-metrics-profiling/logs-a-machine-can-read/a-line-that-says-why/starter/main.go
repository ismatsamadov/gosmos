package main

import (
	"log"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: dropTime}))
	_ = logger
	log.Print("error") // TODO: say what failed and where, with logger.Error
}
