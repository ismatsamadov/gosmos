package main

import "time"

// Line is one JSON log line read back.
type Line struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Msg     string    `json:"msg"`
	Program string    `json:"program"`
	Site    string    `json:"site"`
	PassID  string    `json:"pass_id"`
	Raw     string    `json:"-"` // the line as written
}
