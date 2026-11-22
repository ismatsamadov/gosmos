package main

import "log/slog"

// logResult writes how pass id ended: INFO "pass flown" when err is nil, else ERROR "pass failed".
func logResult(l *slog.Logger, id string, err error) {
	if err != nil {
		l.Error("pass failed", "pass_id", id, "err", err)
		return
	}
	l.Info("pass flown", "pass_id", id)
}
