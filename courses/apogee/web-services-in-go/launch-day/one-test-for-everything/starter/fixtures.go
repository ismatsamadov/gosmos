package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// The customers' test API keys. Production keys never appear in a lesson.
const (
	apogeeKey      = "apogee.9e4d2a7c1b0f3e5d8a6c4b2f0e1d3a5c"
	furrowlightKey = "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"
	brinewatchKey  = "brinewatch.0f1e2d3c4b5a69788796a5b4c3d2e1f0"
	cirrineKey     = "cirrine.c3a1f7e95d2b4086a7e1c9f3b5d8a024"
)

// testTokenKey signs tokens in tests: 32 bytes, never a production key.
const testTokenKey = "bookingd test token key 32 bytes"

// keyHash returns the lowercase hex SHA-256 of an API key, as the service keeps it.
func keyHash(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

// at returns the UTC time an RFC 3339 string names.
func at(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t.UTC()
}

// testStations returns Apogee's five stations, in the bible's order.
func testStations() []Station {
	return []Station{
		{Code: "YST", Name: "Ystmark", LatDeg: 77.8, LonDeg: 2.0, MaskDeg: 5, Bands: []Band{"S", "X"}},
		{Code: "PEL", Name: "Port Ellery", LatDeg: 46.9, LonDeg: -38.1, MaskDeg: 5, Bands: []Band{"UHF", "S"}},
		{Code: "SQT", Name: "Salar Quieto", LatDeg: -23.4, LonDeg: -104.6, MaskDeg: 3, Bands: []Band{"S", "X"}},
		{Code: "LAN", Name: "Lantern Atoll", LatDeg: -1.8, LonDeg: 82.5, MaskDeg: 8, Bands: []Band{"UHF", "S"}},
		{Code: "CWR", Name: "Cape Wren", LatDeg: -53.9, LonDeg: -95.0, MaskDeg: 5, Bands: []Band{"S", "X"}},
	}
}

// testCustomers returns the four customers with their test keys' hashes and webhook secrets, but no webhook URLs.
func testCustomers() []Customer {
	return []Customer{
		{ID: "apogee", Name: "Apogee", SatIDs: []string{"pipit-1", "pipit-2", "pipit-3"}, KeyHash: keyHash(apogeeKey), WebhookKey: []byte("apogee webhook secret")},
		{ID: "furrowlight", Name: "Furrowlight", SatIDs: []string{"furrow-a", "furrow-b"}, KeyHash: keyHash(furrowlightKey), WebhookKey: []byte("furrowlight webhook secret")},
		{ID: "brinewatch", Name: "Brinewatch", SatIDs: []string{"brine-3"}, KeyHash: keyHash(brinewatchKey), WebhookKey: []byte("brinewatch webhook secret")},
		{ID: "cirrine", Name: "Cirrine Atmospherics", SatIDs: []string{"cirrine-1"}, KeyHash: keyHash(cirrineKey), WebhookKey: []byte("cirrine webhook secret")},
	}
}

// testCatalog returns launch week's predicted passes, by AOS then ID.
func testCatalog() []PassView {
	pass := func(id, sat, station, aos, los string, maxElevDeg float64) PassView {
		return PassView{ID: id, Pass: Pass{SatID: sat, Station: station, AOS: at(aos), LOS: at(los), MaxElevDeg: maxElevDeg}}
	}
	return []PassView{
		pass("pipit-1/PEL/09003", "pipit-1", "PEL", "2028-01-31T04:14:26Z", "2028-01-31T04:23:48Z", 65.7),
		pass("furrow-a/SQT/16716", "furrow-a", "SQT", "2028-01-31T14:21:45Z", "2028-01-31T14:31:51Z", 80.1),
		pass("furrow-b/SQT/16714", "furrow-b", "SQT", "2028-01-31T15:09:39Z", "2028-01-31T15:18:44Z", 22.3),
		pass("brine-3/LAN/10088", "brine-3", "LAN", "2028-02-01T12:57:33Z", "2028-02-01T13:03:45Z", 18.3),
		pass("cirrine-1/YST/12218", "cirrine-1", "YST", "2028-02-01T13:30:16Z", "2028-02-01T13:40:20Z", 49.6),
		pass("furrow-a/SQT/16731", "furrow-a", "SQT", "2028-02-01T14:02:10Z", "2028-02-01T14:11:52Z", 47.9),
		pass("furrow-b/SQT/16729", "furrow-b", "SQT", "2028-02-01T14:49:30Z", "2028-02-01T14:58:58Z", 44.0),
		pass("brine-3/LAN/10103", "brine-3", "LAN", "2028-02-02T12:27:50Z", "2028-02-02T12:34:37Z", 22.5),
		pass("cirrine-1/YST/12233", "cirrine-1", "YST", "2028-02-02T13:38:14Z", "2028-02-02T13:48:18Z", 50.4),
	}
}
