package main

import (
	"strconv"
	"strings"
)

// stationReport is the text a host gets on every replan: a heading,
// then a line for each of the station's contacts, in AOS order.
func stationReport(code string, contacts []Contact) string {
	head := code + " contacts: " + strconv.Itoa(len(contacts)) + "\n"
	size := len(head)
	for _, c := range contacts {
		size += len(c.ID) + len(" dish 1\n")
	}
	var b strings.Builder
	b.Grow(size)
	b.WriteString(head)
	for _, c := range contacts {
		b.WriteString(c.ID)
		b.WriteString(" dish ")
		b.WriteString(strconv.Itoa(c.Dish))
		b.WriteByte('\n')
	}
	return b.String()
}
