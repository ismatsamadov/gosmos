package main

import "time"

// Config is passplan 2's settings.
type Config struct {
	Slots      int           `check:"min=1,max=64" default:"12"`
	Days       int           `check:"min=1,max=120" default:"100"`
	MinElevDeg float64       `check:"min=0,max=45" default:"10"`
	Turnaround time.Duration `check:"min=0s,max=10m" default:"2m"`
	Replan     time.Duration `check:"min=1m,max=1h" default:"10m"`
	RunLimit   time.Duration `check:"min=1s,max=5m" default:"30s"`
	Queue      int           `check:"min=1,max=4096" default:"256"`
	Quiet      time.Duration `check:"min=0s,max=1m" default:"5s"`
}
