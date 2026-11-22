package main

// glance is one candidate of take-over day's window, as Tomás's strip on
// the test Board draws it.
type glance struct {
	AOS        int64   // Unix seconds
	LOS        int64   // Unix seconds
	MaxElevDeg float32 // the highest point of the pass
	Sat        uint8   // the satellite's place in fleet order
	Station    uint8   // the station's place in table order
	Dish       uint8   // 1 or 2; 0 when no dish was free
	Booked     bool    // a customer booked the pass
}
