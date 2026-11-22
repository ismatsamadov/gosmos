package main

// Dish is the antenna a site tracks satellites with.
type Dish struct {
	DiameterM float64
	XBand     bool
}

// Station is one of Apogee's ground stations.
type Station struct {
	Code      string
	Equipment Dish
}

// Van is the mobile test van, which carries its own dish.
type Van struct {
	Base      string // code of the station the van works from
	Equipment Dish
}

func main() {
	yst := Station{Code: "YST", Equipment: Dish{DiameterM: 3.7, XBand: true}}
	van := Van{Base: "PEL", Equipment: Dish{DiameterM: 1.2}}
	printDishes(yst, van)
}
