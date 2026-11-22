package main

import "fmt"

const earthRadiusKm = 6371

func main() {
	fmt.Println("Pipit-1 orbit radius:", earthRadiusKm+525, "km")
	fmt.Println("Thistle orbit radius:", earthRadiusKm+415, "km")
	fmt.Println("Cirrine-1 orbit radius:", earthRadiusKm+600, "km")
}
