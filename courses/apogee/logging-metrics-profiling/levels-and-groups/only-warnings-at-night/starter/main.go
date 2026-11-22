package main

import (
	"log/slog"
	"os"
)

func main() {
	opts := &slog.HandlerOptions{ReplaceAttr: dropTime} // TODO: write only WARN and above
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	logger.Info("pass started", "site", "LAN")
	logger.Debug("worker taken", "site", "LAN")
	logger.Debug("lines sent", "site", "LAN")
	logger.Warn("carrier lost", "site", "LAN")
	logger.Error("pass failed", "site", "LAN")
}
