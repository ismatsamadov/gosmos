package main

// Dish is the antenna a site tracks satellites with.
type Dish struct {
	DiameterM float64
	XBand     bool
}

// Station is one of Apogee's ground stations.
type Station struct {
	Code string
	Dish
}

// Van is the mobile test van, which carries its own dish.
type Van struct {
	Base string // code of the station the van works from
	Dish
}

func main() {
	yst := Station{Code: "YST", Dish: Dish{DiameterM: 3.7, XBand: true}}
	van := Van{Base: "PEL", Dish: Dish{DiameterM: 1.2}}
	printDishes(yst, van)
}
