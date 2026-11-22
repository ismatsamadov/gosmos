package main

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
