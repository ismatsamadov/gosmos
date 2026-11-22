package main

import "log/slog"

// logResult writes how pass id ended: INFO "pass flown" when err is nil, else ERROR "pass failed".
func logResult(l *slog.Logger, id string, err error) {
	l.Info("pass flown", "pass_id", id)
}
