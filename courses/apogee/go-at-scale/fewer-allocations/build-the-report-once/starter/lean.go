package main

import "strconv"

// stationReport is the text a host gets on every replan: a heading,
// then a line for each of the station's contacts, in AOS order.
func stationReport(code string, contacts []Contact) string {
	r := code + " contacts: " + strconv.Itoa(len(contacts)) + "\n"
	for _, c := range contacts {
		r += c.ID + " dish " + strconv.Itoa(c.Dish) + "\n"
	}
	return r
}
