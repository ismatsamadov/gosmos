package main

import "fmt"

// Radio is any radio a pass can use.
type Radio interface {
	RateMbps() float64
}

// DopplerTracker is a radio that corrects Doppler shift on its own.
type DopplerTracker interface {
	TracksDoppler()
}

// UHFRadio is the UHF radio on a yagi array.
type UHFRadio struct{}

// RateMbps returns the UHF downlink rate.
func (u UHFRadio) RateMbps() float64 { return 0.0096 }

// SBandRadio is a dish's S-band radio.
type SBandRadio struct{}

// RateMbps returns the S-band downlink rate.
func (s SBandRadio) RateMbps() float64 { return 4 }

// TracksDoppler marks the S-band radio as correcting Doppler shift itself.
func (s SBandRadio) TracksDoppler() {}

// XBandRadio is a dish's X-band radio for imagery.
type XBandRadio struct{}

// RateMbps returns the X-band downlink rate.
func (x XBandRadio) RateMbps() float64 { return 150 }

// TracksDoppler marks the X-band radio as correcting Doppler shift itself.
func (x XBandRadio) TracksDoppler() {}

// needsPrecompensation reports whether the routine must correct r's Doppler shift for it.
func needsPrecompensation(r Radio) bool {
	_, ok := r.(DopplerTracker)
	return !ok
}

func main() {
	fmt.Println("UHF needs precompensation:", needsPrecompensation(UHFRadio{}))
	fmt.Println("S-band needs precompensation:", needsPrecompensation(SBandRadio{}))
}
