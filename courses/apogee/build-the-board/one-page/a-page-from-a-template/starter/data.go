package main

import "time"

// Booking is one reserved pass, as the booking service keeps it.
type Booking struct {
	ID      string    `json:"id"` // "bk-4988", given by the booking service
	SatID   string    `json:"sat_id"`
	Station string    `json:"station"` // station code
	AOS     time.Time `json:"aos"`
	LOS     time.Time `json:"los"`
	Note    string    `json:"note,omitempty"` // written by the customer
}

// Weather is one station's latest reading from the weather service.
type Weather struct {
	Station string    `json:"station"`
	At      time.Time `json:"at"`
	WindKmh float64   `json:"wind_kmh"`
	GustKmh float64   `json:"gust_kmh"`
	TempC   float64   `json:"temp_c"`
}

// Station is one of Apogee's ground stations, with the fields the Board uses.
type Station struct {
	Code string `json:"code"` // "YST": three capital letters
	Name string `json:"name"` // "Ystmark"
}

// stationTable returns Apogee's five ground stations, in the order the Board shows them.
func stationTable() []Station {
	return []Station{
		{Code: "YST", Name: "Ystmark"},
		{Code: "PEL", Name: "Port Ellery"},
		{Code: "SQT", Name: "Salar Quieto"},
		{Code: "LAN", Name: "Lantern Atoll"},
		{Code: "CWR", Name: "Cape Wren"},
	}
}

// rota returns the on-call rota, oldest shift first. Every shift starts on a Monday at
// 08:00 UTC; Tomás took launch week himself.
func rota() []Shift {
	monday := func(month time.Month, day int) time.Time {
		return time.Date(2028, month, day, 8, 0, 0, 0, time.UTC)
	}
	return []Shift{
		{From: monday(time.January, 31), Name: "Tomás Ferreira-Lund"},
		{From: monday(time.February, 7), Name: "Nkem Halloran"},
		{From: monday(time.February, 14), Name: "Mirela Dunn"},
	}
}
