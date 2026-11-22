package main

import "fmt"

// DishState is what a dish is doing.
type DishState uint8

// The states a controller reports.
const (
	Idle DishState = iota
	Moving
	Tracking
	Stowed
	Fault
)

// String names the state as the controller does: "TRACKING".
func (s DishState) String() string {
	switch s {
	case Idle:
		return "IDLE"
	case Moving:
		return "MOVING"
	case Tracking:
		return "TRACKING"
	case Stowed:
		return "STOWED"
	case Fault:
		return "FAULT"
	}
	return fmt.Sprintf("STATE(%d)", uint8(s))
}

// The status frame's flags.
const (
	FlagHeater   uint16 = 1 << 0
	FlagWindLock uint16 = 1 << 1
	FlagFault    uint16 = 1 << 2
)

// Status is a decoded status frame: one dish at one moment.
type Status struct {
	Dish    uint8
	State   DishState
	AzDeg   float64 // 183.25
	ElDeg   float64 // 41.70
	WindKmh float64 // 64.0
	TempC   float64 // -12.5
	Flags   uint16
	LastSeq uint32
}
