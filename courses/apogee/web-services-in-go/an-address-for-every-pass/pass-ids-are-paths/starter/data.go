package main

import (
	"strings"
	"time"
)

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

// catalog returns passplan's predicted passes from 25 November to 10 December 2027, an excerpt, ordered by AOS, then ID.
func catalog() []PassView {
	return []PassView{
		predicted("pipit-1/PEL/07987", "2027-11-25T03:49:57Z", "2027-11-25T03:58:52Z", 33.3),
		predicted("furrow-a/SQT/15697", "2027-11-25T14:31:40Z", "2027-11-25T14:41:39Z", 54.4),
		predicted("brine-3/LAN/09056", "2027-11-26T02:24:30Z", "2027-11-26T02:31:13Z", 21.9),
		predicted("furrow-b/CWR/15709", "2027-11-26T13:33:23Z", "2027-11-26T13:42:37Z", 73.5),
		predicted("cirrine-1/YST/11260", "2027-11-29T08:09:40Z", "2027-11-29T08:19:55Z", 59.3),
		predicted("furrow-a/SQT/15758", "2027-11-29T14:46:40Z", "2027-11-29T14:56:15Z", 31.5),
		predicted("pipit-1/PEL/08078", "2027-12-01T03:53:31Z", "2027-12-01T04:02:32Z", 36.6),
		predicted("furrow-b/SQT/15786", "2027-12-01T14:54:30Z", "2027-12-01T15:04:15Z", 36.7),
		predicted("brine-3/LAN/09168", "2027-12-03T10:42:22Z", "2027-12-03T10:50:29Z", 83.9),
		predicted("furrow-a/SQT/15818", "2027-12-03T13:28:42Z", "2027-12-03T13:36:51Z", 15.3),
		predicted("pipit-1/PEL/08154", "2027-12-06T04:12:05Z", "2027-12-06T04:21:27Z", 61.4),
		predicted("cirrine-1/YST/11364", "2027-12-06T07:30:00Z", "2027-12-06T07:40:04Z", 44.2),
		predicted("brine-3/LAN/09213", "2027-12-06T09:14:20Z", "2027-12-06T09:22:02Z", 38.5),
		predicted("furrow-a/SQT/15864", "2027-12-06T14:02:10Z", "2027-12-06T14:11:52Z", 47.9),
		predicted("furrow-b/SQT/15862", "2027-12-06T14:49:30Z", "2027-12-06T14:58:58Z", 44.0),
		predicted("furrow-a/CWR/15879", "2027-12-07T13:50:25Z", "2027-12-07T13:59:15Z", 35.9),
		predicted("furrow-b/YST/15890", "2027-12-08T10:33:15Z", "2027-12-08T10:42:28Z", 75.3),
		predicted("furrow-a/SQT/15910", "2027-12-09T14:36:38Z", "2027-12-09T14:46:32Z", 44.9),
		predicted("cirrine-1/YST/11424", "2027-12-10T08:02:06Z", "2027-12-10T08:12:20Z", 56.1),
		predicted("furrow-b/SQT/15923", "2027-12-10T15:04:35Z", "2027-12-10T15:13:55Z", 26.2),
	}
}

// predicted returns the catalog entry for the pass with this ID, AOS, LOS and maximum elevation.
func predicted(id, aos, los string, maxElevDeg float64) PassView {
	parts := strings.Split(id, "/") // satellite, station, orbit
	return PassView{ID: id, Pass: Pass{SatID: parts[0], Station: parts[1], AOS: mustTime(aos), LOS: mustTime(los), MaxElevDeg: maxElevDeg}}
}

// mustTime parses one of this file's fixed RFC 3339 times; it panics only on a typo here.
func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
