package main

import "fmt"

// Radio is any radio a pass can use.
type Radio interface {
	RateMbps() float64
}

// XBandRadio is a dish's X-band radio for imagery.
type XBandRadio struct{}

// RateMbps returns the X-band downlink rate.
func (x XBandRadio) RateMbps() float64 { return 150 }

// KaBandRadio is the Ka-band test radio on La Grande at Salar Quieto.
type KaBandRadio struct{}

// RateMbps returns the Ka-band downlink rate.
func (k KaBandRadio) RateMbps() float64 { return 300 }

// downlinkSeconds returns how many seconds r needs to bring megabytes down.
func downlinkSeconds(r Radio, megabytes float64) float64 {
	return megabytes * 8 / r.RateMbps()
}

func main() {
	fmt.Printf("Ka-band: 9000 MB in %.0f s\n", downlinkSeconds(KaBandRadio{}, 9000))
}
