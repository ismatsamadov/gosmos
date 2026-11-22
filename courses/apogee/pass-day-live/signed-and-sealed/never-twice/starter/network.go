package main

import (
	"fmt"
	"time"
)

// Site is a station as passplan sees it.
type Site struct {
	Station
	Dishes int       // the dishes passplan can give passes to
	Live   time.Time // the first moment it flies a pass; zero means always
}

// Elements place a satellite on its orbit for the pass model.
type Elements struct {
	SatID    string
	Orbit    Orbit     // altitude and inclination
	Epoch    time.Time // the moment the numbers below describe; no pass before it
	Rev0     int       // the orbit number at Epoch
	NodeDeg  float64   // east longitude of the northbound equator crossing, at Epoch
	PhaseDeg float64   // how far past that crossing the satellite is at Epoch, 0 to 360
	// RevsPerDay is the satellite's orbits a day, where tracking measured them; zero means
	// the period its altitude gives.
	RevsPerDay float64
}

// Network is everything passplan plans for.
type Network struct {
	Sites      []Site      // the twenty, in table order
	Satellites []Satellite // the eighty, in fleet order
	Elements   []Elements  // the same eighty, in the same order
}

// newNetwork returns the network as it stands at the end of June 2029.
func newNetwork() *Network {
	sx, s, us := []Band{SBand, XBand}, []Band{SBand}, []Band{UHF, SBand}
	w1 := time.Date(2029, 4, 2, 0, 0, 0, 0, time.UTC)
	w2 := time.Date(2029, 5, 7, 0, 0, 0, 0, time.UTC)
	w3 := time.Date(2029, 6, 4, 0, 0, 0, 0, time.UTC)
	site := func(code, name string, lat, lon, mask float64, bands []Band, dishes int, live time.Time) Site {
		st := Station{Code: code, Name: name, LatDeg: lat, LonDeg: lon, MaskDeg: mask, Bands: bands}
		return Site{Station: st, Dishes: dishes, Live: live}
	}
	var always time.Time
	n := &Network{Sites: []Site{
		site("YST", "Ystmark", 77.8, 2.0, 5, sx, 2, always),
		site("PEL", "Port Ellery", 46.9, -38.1, 5, us, 1, always),
		site("SQT", "Salar Quieto", -23.4, -104.6, 3, sx, 2, always),
		site("LAN", "Lantern Atoll", -1.8, 82.5, 8, us, 1, always),
		site("CWR", "Cape Wren", -53.9, -95.0, 5, sx, 1, always),
		site("TIR", "Tirrin Rock", 69.5, -12.0, 5, sx, 2, w1),
		site("VSK", "Varsk Island", 51.0, -160.0, 7, sx, 1, w1),
		site("CAL", "Calidra Atoll", 0.5, -150.0, 8, s, 1, w1),
		site("BRM", "Ilha Brumosa", -30.0, -15.0, 5, sx, 1, w1),
		site("OST", "Ostrel Island", -42.5, 70.0, 5, sx, 1, w1),
		site("ILV", "Ilvane Spit", 72.0, -165.0, 3, sx, 2, w2),
		site("CSL", "Cape Solvane", 35.0, -45.0, 5, sx, 1, w2),
		site("SRH", "Saltreach Island", 17.0, 135.0, 6, s, 1, w2),
		site("CRP", "Corith Point", 10.0, 65.0, 6, s, 1, w2),
		site("NIM", "Nimrie Point", -45.0, -150.0, 5, sx, 1, w2),
		site("HSK", "Hasku Point", 32.0, 170.0, 4, sx, 1, w3),
		site("ICM", "Ilha Calmaria", 8.0, -30.0, 7, s, 1, w3),
		site("QND", "Quindle Island", -25.0, 95.0, 5, sx, 1, w3),
		site("PVT", "Punta Ventisca", -60.5, 40.0, 5, sx, 1, w3),
		site("VRL", "Vorrel Island", -64.0, -140.0, 4, sx, 2, w3),
	}}
	add := func(id, name, op string, o Orbit, bands []Band, epoch time.Time, rev0 int, node, phase, revs float64) {
		n.Satellites = append(n.Satellites, Satellite{ID: id, Name: name, Operator: op, Orbit: o, Bands: bands, Mode: "NOMINAL"})
		n.Elements = append(n.Elements, Elements{SatID: id, Orbit: o, Epoch: epoch, Rev0: rev0, NodeDeg: node, PhaseDeg: phase, RevsPerDay: revs})
	}
	jan := time.Date(2029, 1, 1, 0, 0, 0, 0, time.UTC)
	xs := []Band{XBand, SBand}
	// Pipit-1, Pipit-2, Furrow-A, Furrow-B, Brine-3 and Cirrine-1 fitted to passes the
	// stations recorded (Furrow-A and Furrow-B make 867 orbits in 57 days); Pipit-3 from
	// Apogee's own tracking.
	add("pipit-1", "Pipit-1", "Apogee", Orbit{525, 97.5}, []Band{UHF}, jan, 14094, 201.4926025817, 176.9565679403, 0)
	add("pipit-2", "Pipit-2", "Apogee", Orbit{550, 97.6}, s, jan, 19752, 234.8439482122, 339.4065995135, 0)
	add("pipit-3", "Pipit-3", "Apogee", Orbit{510, 97.4}, sx, jan, 9689, 150, 40, 0)
	add("furrow-a", "Furrow-A", "Furrowlight", Orbit{510, 97.4}, xs, jan, 21818, 294.5055051072, 52.6136621172, 867.0/57)
	add("furrow-b", "Furrow-B", "Furrowlight", Orbit{510, 97.4}, xs, jan, 21815, 297.5851231763, 231.0891932221, 867.0/57)
	add("brine-3", "Brine-3", "Brinewatch", Orbit{480, 53.0}, s, jan, 15209, 179.2270109042, 188.3141532885, 0)
	add("cirrine-1", "Cirrine-1", "Cirrine Atmospherics", Orbit{600, 97.8}, s, jan, 17206, 282.5236543824, 268.6744811484, 0)
	add("thistle-2", "Thistle-2", "Kessick Polytechnic", Orbit{505, 97.4}, []Band{UHF},
		time.Date(2029, 3, 8, 10, 20, 0, 0, time.UTC), 1, -120, 0, 0)
	for k := 2; k <= 13; k++ {
		add(fmt.Sprintf("cirrine-%d", k), fmt.Sprintf("Cirrine-%d", k), "Cirrine Atmospherics", Orbit{600, 97.8}, s,
			time.Date(2029, 4, 17, 9, 40, 0, 0, time.UTC), 1, -150, float64(30*(k-2)), 0)
	}
	for k := 4; k <= 39; k++ {
		add(fmt.Sprintf("brine-%d", k), fmt.Sprintf("Brine-%d", k), "Brinewatch", Orbit{480, 53.0}, s,
			time.Date(2029, 5, 16, 13, 30, 0, 0, time.UTC), 1, 120, float64(10*(k-4)), 0)
	}
	for k := 0; k < 24; k++ {
		add(fmt.Sprintf("furrow-%c", 'c'+k), fmt.Sprintf("Furrow-%c", 'C'+k), "Furrowlight", Orbit{510, 97.4}, xs,
			time.Date(2029, 6, 27, 20, 40, 0, 0, time.UTC), 1, 10, float64(15*k), 0)
	}
	return n
}
