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

// HostConfig is the form a station's host fills in when its dish is installed.
type HostConfig struct {
	Code    string   `check:"required,len=3"`
	Name    string   `check:"required"`
	LatDeg  float64  `check:"min=-90,max=90"`
	LonDeg  float64  `check:"min=-180,max=180"`
	MaskDeg float64  `check:"min=0,max=15"`
	Dishes  int      `check:"min=1,max=2"`
	Bands   []string `check:"required,oneof=UHF S X"`
}

// PipelineConfig is tlmd's ingest settings.
type PipelineConfig struct {
	Workers int      `check:"min=1,max=32" default:"4"`
	Buffer  int      `check:"min=0,max=1024" default:"16"`
	Feeds   []string `check:"required,len=3"` // station codes
}
