package main

import "iter"

// At yields the dish and the contact of each of the station's
// contacts, in the schedule's order.
func (sch *Schedule) At(station string) iter.Seq2[int, Contact] {
	return func(yield func(int, Contact) bool) {
		for _, c := range sch.Contacts {
			if c.Station == station && !yield(c.Dish, c) {
				return
			}
		}
	}
}
