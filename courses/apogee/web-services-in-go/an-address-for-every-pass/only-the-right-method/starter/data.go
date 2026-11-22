package main

// stationTable returns Apogee's five ground stations in the table's order, with the bands a customer can book.
func stationTable() []Station {
	return []Station{
		{Code: "YST", Name: "Ystmark", LatDeg: 77.8, LonDeg: 2.0, MaskDeg: 5, Bands: []Band{SBand, XBand}},
		{Code: "PEL", Name: "Port Ellery", LatDeg: 46.9, LonDeg: -38.1, MaskDeg: 5, Bands: []Band{UHF, SBand}},
		{Code: "SQT", Name: "Salar Quieto", LatDeg: -23.4, LonDeg: -104.6, MaskDeg: 3, Bands: []Band{SBand, XBand}},
		{Code: "LAN", Name: "Lantern Atoll", LatDeg: -1.8, LonDeg: 82.5, MaskDeg: 8, Bands: []Band{UHF, SBand}},
		{Code: "CWR", Name: "Cape Wren", LatDeg: -53.9, LonDeg: -95.0, MaskDeg: 5, Bands: []Band{SBand, XBand}},
	}
}
