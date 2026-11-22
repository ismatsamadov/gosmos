package main

import "fmt"

// UHFRadio is the UHF radio on a yagi array.
type UHFRadio struct{}

// RateMbps returns the UHF downlink rate.
func (u UHFRadio) RateMbps() float64 { return 0.0096 }

// SBandRadio is a dish's S-band radio.
type SBandRadio struct{}

// RateMbps returns the S-band downlink rate.
func (s SBandRadio) RateMbps() float64 { return 4 }

// XBandRadio is a dish's X-band radio for imagery.
type XBandRadio struct{}

// RateMbps returns the X-band downlink rate.
func (x XBandRadio) RateMbps() float64 { return 150 }

// Radio is any radio a pass can use.
type Radio interface {
	RateMbps() float64
}

// downlinkSeconds returns how many seconds r needs to bring megabytes down.
func downlinkSeconds(r Radio, megabytes float64) float64 {
	return megabytes * 8 / r.RateMbps()
}

func main() {
	fmt.Printf("X-band: 9000 MB in %.0f s\n", downlinkSeconds(XBandRadio{}, 9000))
}
